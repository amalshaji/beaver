package main

import (
	"fmt"
	"strconv"

	"github.com/amalshaji/beaver/client"
	"github.com/spf13/cobra"
)

var (
	port      int
	subdomain string
	httpCmd   = &cobra.Command{
		Use:   "http [PORT]",
		Short: "Tunnel local http servers",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("local server port is required")
			}
			if len(args) > 1 {
				return fmt.Errorf("only one port number is allowed")
			}

			var err error
			port, err = strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("port must be a number")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var tunnels = make([]client.TunnelConfig, 0)
			tunnels = append(tunnels, client.TunnelConfig{Port: port, Subdomain: subdomain})
			startTunnels(tunnels)
		},
	}
)

func init() {
	httpCmd.Flags().StringVar(&subdomain, "subdomain", "", "Subdomain to tunnel http requests (default \"<random_subdomain>\")")

	rootCmd.AddCommand(httpCmd)
}
