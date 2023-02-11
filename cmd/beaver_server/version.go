package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version number of beaver server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Beaver server v0.0.2")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
