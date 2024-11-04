/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"awesome-aws-cli/cmd/utils" // Replace with your actual module path

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

// s3Cmd represents the s3 command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Manage S3 buckets and objects",
	Run: func(cmd *cobra.Command, args []string) {
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")

		_, err := utils.LoadAWSConfig(profile, region)
		if err != nil {
			fmt.Printf("Error loading AWS config: %v\n", err)
			return
		}

		// ... (rest of the function)
	},
}

func init() {
	rootCmd.AddCommand(s3Cmd)

	// Add S3 subcommands
	s3Cmd.AddCommand(listS3BucketsCmd())
	// Add more S3 subcommands here
}

func listS3BucketsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-buckets",
		Short: "List S3 buckets",
		Run: func(cmd *cobra.Command, args []string) {
			profile, _ := cmd.Flags().GetString("profile")
			region, _ := cmd.Flags().GetString("region")
			cfg, err := utils.LoadAWSConfig(profile, region)
			if err != nil {
				fmt.Printf("Error loading AWS config: %v\n", err)
				return
			}

			client := s3.NewFromConfig(cfg)

			resp, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
			if err != nil {
				fmt.Printf("Unable to list buckets: %v\n", err)
				return
			}

			for _, bucket := range resp.Buckets {
				fmt.Printf("Bucket Name: %s, Creation Date: %s\n", *bucket.Name, bucket.CreationDate)
			}
		},
	}
}

// Implement more S3 commands here (e.g., create bucket, delete bucket, list objects, etc.)
