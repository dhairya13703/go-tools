/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"awesome-aws-cli/cmd/utils" // Replace with your actual module path

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/spf13/cobra"
)

// lambdaCmd represents the lambda command
var lambdaCmd = &cobra.Command{
	Use:   "lambda",
	Short: "Manage AWS Lambda functions",
	Run: func(cmd *cobra.Command, args []string) {
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")

		cfg, err := utils.LoadAWSConfig(profile, region)
		if err != nil {
			fmt.Printf("Error loading AWS config: %v\n", err)
			return
		}

		client := lambda.NewFromConfig(cfg)

		resp, err := client.ListFunctions(context.TODO(), &lambda.ListFunctionsInput{})
		if err != nil {
			fmt.Printf("Unable to list functions: %v\n", err)
			return
		}

		for _, function := range resp.Functions {
			fmt.Printf("Function Name: %s, Runtime: %s, Last Modified: %s\n", *function.FunctionName, function.Runtime, *function.LastModified)
		}
	},
}

func init() {
	rootCmd.AddCommand(lambdaCmd)

	// Add Lambda subcommands
	lambdaCmd.AddCommand(listLambdaFunctionsCmd())
	// Add more Lambda subcommands here
}

func listLambdaFunctionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List Lambda functions",
		Run: func(cmd *cobra.Command, args []string) {
			profile, _ := cmd.Flags().GetString("profile")
			region, _ := cmd.Flags().GetString("region")
			cfg, err := utils.LoadAWSConfig(profile, region)
			if err != nil {
				fmt.Printf("Error loading AWS config: %v\n", err)
				return
			}

			client := lambda.NewFromConfig(cfg)

			resp, err := client.ListFunctions(context.TODO(), &lambda.ListFunctionsInput{})
			if err != nil {
				fmt.Printf("Unable to list functions: %v\n", err)
				return
			}

			for _, function := range resp.Functions {
				fmt.Printf("Function Name: %s, Runtime: %s, Last Modified: %s\n", *function.FunctionName, function.Runtime, *function.LastModified)
			}
		},
	}
}

// Implement more Lambda commands here (e.g., invoke function, update function, etc.)
