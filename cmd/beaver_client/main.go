package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/amalshaji/beaver/client"
)

func main() {
	ctx := context.Background()

	configFile := flag.String("config", "", "config file path")
	subdomain := flag.String("subdomain", "", "subdomain to create the tunnel at")
	port := flag.Int("port", 0, "local server port to tunnel")

	flag.Parse()

	if *configFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		*configFile = fmt.Sprintf("%s/.beaver/beaver_client.yaml", homeDir)
	}

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
