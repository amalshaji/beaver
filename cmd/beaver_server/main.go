package main

import (
	"log"

	handler "github.com/amalshaji/beaver/internal/server/handlers"
)

func startServer() {
	handler.Start(configFile)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
