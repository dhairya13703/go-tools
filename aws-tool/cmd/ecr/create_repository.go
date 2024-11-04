package ecr

import (
	"context"
	"fmt"

	"awesome-aws-cli/cmd/utils"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/spf13/cobra"
)

func NewCreateRepositoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-repository [repository-name]",
		Short: "Create a new ECR repository",
		Args:  cobra.ExactArgs(1),
		Run:   runCreateRepository,
	}

	return cmd
}

func runCreateRepository(cmd *cobra.Command, args []string) {
	repositoryName := args[0]
	profile, _ := cmd.Flags().GetString("profile")
	region, _ := cmd.Flags().GetString("region")

	cfg, err := utils.LoadAWSConfig(profile, region)
	if err != nil {
		fmt.Printf("Error loading AWS config: %v\n", err)
		return
	}

	client := ecr.NewFromConfig(cfg)

	input := &ecr.CreateRepositoryInput{
		RepositoryName: &repositoryName,
	}

	result, err := client.CreateRepository(context.TODO(), input)
	if err != nil {
		fmt.Printf("Error creating repository: %v\n", err)
		return
	}

	fmt.Printf("Repository created successfully: %s\n", *result.Repository.RepositoryUri)
}
