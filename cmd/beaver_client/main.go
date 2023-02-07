package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/amalshaji/beaver/client"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

func getDefaultConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/.beaver/beaver_client.yaml", homeDir)
}

func loadProxyConfig(path string) (*client.Proxy, error) {
	var config client.Proxy

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	ctx := context.Background()

	configFile := flag.String("config", getDefaultConfigFilePath(), "Config file path")
	subdomain := flag.String("subdomain", "", "Subdomain to tunnel http requests (default \"<random_subdomain>\")")
	port := flag.Int("port", 0, "Local http server port (required)")
	startAll := flag.Bool("start-all", false, "Start all tunnels defined in config file")
	showWsReadErrors := flag.Bool("showtunnelreaderrors", false, "Enable websocket read errors")

	flag.CommandLine.MarkHidden("showtunnelreaderrors")

	flag.CommandLine.SortFlags = false
	flag.ErrHelp = fmt.Errorf("")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "beaver - tunnel local ports to public URLs:\n\nUsage:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	var proxies []*client.Client

	proxyConfig, err := loadProxyConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	if *startAll {
		if len(proxyConfig.Tunnels) == 0 {
			log.Fatal("No tunnels defined in the config file")
		}
	} else {
		if *port == 0 {
			log.Fatalln("local server port is required")
		}
		if len(proxyConfig.Tunnels) != 0 {
			proxyConfig.Tunnels = []client.TunnelConfig{{Subdomain: *subdomain, Port: *port}}
		}
	}

	for _, tunnel := range proxyConfig.Tunnels {
		config, err := client.LoadConfiguration(*proxyConfig, tunnel.Subdomain, tunnel.Port, *showWsReadErrors)
		if err != nil {
			log.Fatalf("Unable to load configuration : %s", err)
		}
		proxy := client.NewClient(&config)
		proxies = append(proxies, proxy)
		proxy.Start(ctx)
	}

	// Wait signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	// When receives the signal, shutdown
	for _, proxy := range proxies {
		proxy.Shutdown()
	}

}
