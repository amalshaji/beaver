package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/amalshaji/beaver/internal/utils"
	"github.com/labstack/gommon/color"
	"github.com/spf13/cobra"
)

func getDefaultConfigFileDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/.beaver", homeDir)
}

var (
	initConfig bool
	configCmd  = &cobra.Command{
		Use:   "config",
		Short: "Manage client config",
		Run: func(cmd *cobra.Command, args []string) {
			err := createConfigFile()
			if err != nil {
				fmt.Println(color.Red(err.Error()))
			}
		},
	}
	ErrUnableToCreateConfigFile = errors.New("unable to create the config file")
)

var ConfigTemplate string = `
target: 
secretkey:
tunnels:
  - name: tunnel-1
    subdomain: subdomain-1
    port: 8000
`

func createConfigFile() error {
	if initConfig {
		_, err := os.Stat(configFile)
		if !os.IsNotExist(err) {
			fmt.Println(color.Yellow("Client config exists at: " + getDefaultConfigFilePath()))
			return nil
		}
		_, err = os.Stat(getDefaultConfigFileDir())
		if os.IsNotExist(err) {
			err = os.Mkdir(getDefaultConfigFileDir(), os.ModePerm)
			if err != nil {
				return ErrUnableToCreateConfigFile
			}
		}
		f, err := os.Create(getDefaultConfigFilePath())
		if err != nil {
			return ErrUnableToCreateConfigFile
		}
		defer f.Close()
		f.WriteString(utils.SanitizeString(ConfigTemplate))
		fmt.Println(color.Green("Client config created at: " + getDefaultConfigFilePath()))
	}
	return nil
}

func init() {
	configCmd.Flags().BoolVar(&initConfig, "init", false, "Create the default config file template")

	rootCmd.AddCommand(configCmd)
}
