/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package net

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var urlPath string
var client = http.Client{
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 2 * time.Second,
		}).Dial,
	},
}

func ping(domain string) (int, error) {
	url := "https://" + domain
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	res.Body.Close()
	return res.StatusCode, nil
}

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping is a tool that pings a host",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Logic
		if res, err := ping(urlPath); err != nil {
			fmt.Println("Error pinging host:", err)
		} else {
			fmt.Println("Host is reachable with status code:", res)
		}
	},
}

func init() {
	pingCmd.Flags().StringVarP(&urlPath, "url", "u", "", "The url to ping")

	if err := pingCmd.MarkFlagRequired("url"); err != nil {
		fmt.Println("Error marking flag as required:", err)
	}

	NetCmd.AddCommand(pingCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
