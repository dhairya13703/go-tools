package ecr

import (
	"context"
	"fmt"

	"awesome-aws-cli/cmd/utils"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/spf13/cobra"
)

func NewListRepositoriesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-repositories",
		Short: "List ECR repositories",
		Run:   runListRepositories,
	}

	return cmd
}

func runListRepositories(cmd *cobra.Command, args []string) {
	profile, _ := cmd.Flags().GetString("profile")
	region, _ := cmd.Flags().GetString("region")

	cfg, err := utils.LoadAWSConfig(profile, region)
	if err != nil {
		fmt.Printf("Error loading AWS config: %v\n", err)
		return
	}

	client := ecr.NewFromConfig(cfg)

	input := &ecr.DescribeRepositoriesInput{}
	resp, err := client.DescribeRepositories(context.TODO(), input)
	if err != nil {
		fmt.Printf("Error describing repositories: %v\n", err)
		return
	}

	for _, repo := range resp.Repositories {
		fmt.Printf("Repository: %s\n", *repo.RepositoryName)
		fmt.Printf("  URI: %s\n", *repo.RepositoryUri)
		fmt.Printf("  Created: %s\n", repo.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}
}
