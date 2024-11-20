package cmd

import (
	"github.com/spf13/cobra"
)

var (
	configFile string
	rootCmd    = &cobra.Command{
		Use:   "ssh-tool",
		Short: "A tool for managing SSH connections to servers in local machine with ssh connections",
		Long: `A CLI tool that helps manage and connect to various servers 
               using embedded configuration with optional external config file support.`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "optional external config file")
}
