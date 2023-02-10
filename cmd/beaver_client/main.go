package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/amalshaji/beaver/client"
	"gopkg.in/yaml.v3"
)

func loadProxyConfig() (*client.Config, error) {
	var config *client.Config

	bytes, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func startTunnels(tunnels []client.TunnelConfig) {
	ctx := context.Background()
	var proxies []*client.Client

	proxyConfig, err := loadProxyConfig()
	if err != nil {
		log.Fatal(err)
	}

	for _, proxyTunnel := range tunnels {
		config, err := client.LoadConfiguration(*proxyConfig, proxyTunnel.Subdomain, proxyTunnel.Port, showWsReadErrors)
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

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
