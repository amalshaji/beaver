package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "v0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version number of beaver server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Beaver server %s\n", VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
