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
)

func getDefaultConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/.beaver/beaver_client.yaml", homeDir)
}

func main() {
	ctx := context.Background()

	configFile := flag.String("config", getDefaultConfigFilePath(), "Config file path")
	subdomain := flag.String("subdomain", "", "Subdomain to tunnel http requests (default \"<random_subdomain>\")")
	port := flag.Int("port", 0, "Local http server port (required)")

	flag.CommandLine.SortFlags = false
	flag.ErrHelp = fmt.Errorf("")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "beaver - tunnel local ports to public URLs:\n\nUsage:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *port == 0 {
		log.Fatalln("local server port is required")
	}

	// Load configuration
	config, err := client.LoadConfiguration(*configFile, *subdomain, *port)
	if err != nil {
		log.Fatalf("Unable to load configuration : %s", err)
	}

	proxy := client.NewClient(config)
	proxy.Start(ctx)

	// Wait signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	// When receives the signal, shutdown
	proxy.Shutdown()
}
