package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/amalshaji/beaver/server"
)

func startServer() {
	// Load configuration
	config, err := server.LoadConfiguration(configFile)
	if err != nil {
		log.Fatalf("Unable to load configuration : %s", err)
	}

	server := server.NewServer(config)
	server.Start()

	// Wait signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	// When receives the signal, shutdown
	server.Shutdown()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
