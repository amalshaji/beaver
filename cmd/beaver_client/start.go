package main

import (
	"fmt"
	"log"
	"os"

	"github.com/amalshaji/beaver/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	all      bool
	startCmd = &cobra.Command{
		Use:   "start [--all] or [tunnel1 tunnel2]",
		Short: "Start tunnels defined in the config file",
		Args: func(cmd *cobra.Command, args []string) error {
			if !all && len(args) == 0 {
				return fmt.Errorf("either --all or a list of tunnel service names must be passed")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			proxyTunnels, err := loadProxyTunnelConfig()
			if err != nil {
				log.Fatal(err)
			}

			var filteredTunnels = make([]client.TunnelConfig, 0)

			if all {
				filteredTunnels = proxyTunnels.Tunnels
			} else {
				tunnelsToStart := os.Args[2:]
				for _, tunnel := range tunnelsToStart {
					for _, proxyTunnel := range proxyTunnels.Tunnels {
						if proxyTunnel.Name == tunnel {
							filteredTunnels = append(filteredTunnels, proxyTunnel)
						}
					}
				}
			}

			startTunnels(filteredTunnels)
		},
	}
)

func loadProxyTunnelConfig() (*client.ProxyTunnels, error) {
	var config *client.ProxyTunnels

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

func init() {
	startCmd.Flags().BoolVar(&all, "all", false, "Start all tunnels listed in the config")

	rootCmd.AddCommand(startCmd)
}
