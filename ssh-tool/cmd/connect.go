package cmd

import (
	"fmt"
	"ssh-tool/internal/config"
	"ssh-tool/internal/ssh"
	"strconv"

	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect [server number]",
	Short: "Connect to a server by its number",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		num, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Invalid server number: %v\n", err)
			return
		}

		servers := cfg.GetServersList()
		if num < 1 || num > len(servers) {
			fmt.Printf("Invalid server number. Please choose between 1 and %d\n", len(servers))
			return
		}

		server := servers[num-1]
		fmt.Printf("Connecting to %s (%s)...\n", server.Name, server.Hostname)

		client := ssh.NewClient(server)
		if err := client.Connect(); err != nil {
			fmt.Printf("Error connecting to server: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
