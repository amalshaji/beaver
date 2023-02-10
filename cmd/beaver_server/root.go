package main

import "github.com/spf13/cobra"

var (
	configFile string
	rootCmd    = &cobra.Command{
		Use:   "beaver_server",
		Short: "Tunnel local ports to public URLs",
		Run: func(cmd *cobra.Command, args []string) {
			startServer()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "beaver_server.yaml", "Path to config file")
}
