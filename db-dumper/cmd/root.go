package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	hostname string
	username string
	password string
	port     string
	database string
	output   string
)

var rootCmd = &cobra.Command{
	Use:   "db-dumper",
	Short: "A database backup utility",
	Long:  `db-dumper is a CLI utility to backup different types of databases`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&hostname, "host", "localhost", "Database host")
	rootCmd.PersistentFlags().StringVar(&username, "user", "", "Database username")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Database password")
	rootCmd.PersistentFlags().StringVar(&port, "port", "", "Database port")
	rootCmd.PersistentFlags().StringVar(&database, "db", "", "Database name")
	rootCmd.PersistentFlags().StringVar(&output, "output", "", "Output file path")
}
