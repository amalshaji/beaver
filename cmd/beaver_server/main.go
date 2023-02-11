package main

import (
	"log"

	"github.com/amalshaji/beaver/internal/server"
)

func startServer() {
	server.Start(configFile)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
