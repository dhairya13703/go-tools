/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"awesome-aws-cli/cmd/cloudwatch"
	"awesome-aws-cli/cmd/ec2"
	"awesome-aws-cli/cmd/ecr"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "awscli",
	Short: "An amazing AWS CLI for DevOps engineers",
	Long: `A custom AWS CLI tool with enhanced features for efficient DevOps workflows.
This CLI provides improved options and efficiency for daily AWS operations.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("profile", "p", "", "AWS profile to use")
	rootCmd.PersistentFlags().StringP("region", "r", "", "AWS region to use")

	viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))
	viper.BindPFlag("region", rootCmd.PersistentFlags().Lookup("region"))

	rootCmd.AddCommand(ec2.NewEC2Cmd())
	rootCmd.AddCommand(cloudwatch.NewCloudWatchCmd())
	rootCmd.AddCommand(ecr.NewECRCmd())
	// Add other top-level commands here
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".awscli")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
