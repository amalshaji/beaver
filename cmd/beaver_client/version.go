package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "0.3.0-alpha.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print beaver client version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Beaver client %s\n", VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
