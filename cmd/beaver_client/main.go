package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/amalshaji/beaver/internal/client"
)

func startTunnels(tunnels []client.TunnelConfig) {
	ctx := context.Background()
	var proxies []*client.Client

	for _, proxyTunnel := range tunnels {
		config, err := client.LoadConfiguration(configFile, proxyTunnel.Subdomain, proxyTunnel.Port, showWsReadErrors)
		if err != nil {
			log.Fatalf("Unable to load configuration: %s", err)
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
