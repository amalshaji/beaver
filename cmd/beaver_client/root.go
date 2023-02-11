package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configFile       string
	showWsReadErrors bool
	rootCmd          = &cobra.Command{
		Use:   "beaver",
		Short: "Tunnel local ports to public URLs",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func getDefaultConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/.beaver/beaver_client.yaml", homeDir)
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", getDefaultConfigFilePath(), "Path to the client config file")
	rootCmd.PersistentFlags().BoolVar(&showWsReadErrors, "showWsReadErrors", false, "Log websocket read errors")

	rootCmd.PersistentFlags().MarkHidden("showWsReadErrors")
}
