package ecr

import (
	"fmt"
	"os/exec"

	"awesome-aws-cli/cmd/utils"

	"github.com/spf13/cobra"
)

func NewPullImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull [repository-name] [image-tag]",
		Short: "Pull an image from ECR",
		Args:  cobra.ExactArgs(2),
		Run:   runPullImage,
	}

	return cmd
}

func runPullImage(cmd *cobra.Command, args []string) {
	repositoryName := args[0]
	imageTag := args[1]
	profile, _ := cmd.Flags().GetString("profile")
	region, _ := cmd.Flags().GetString("region")

	cfg, err := utils.LoadAWSConfig(profile, region)
	if err != nil {
		fmt.Printf("Error loading AWS config: %v\n", err)
		return
	}

	// Get the ECR registry URL
	registryURL, err := utils.GetECRRegistryURL(cfg)
	if err != nil {
		fmt.Printf("Error getting ECR registry URL: %v\n", err)
		return
	}

	// Construct the full image name
	fullImageName := fmt.Sprintf("%s/%s:%s", registryURL, repositoryName, imageTag)

	// Pull the image
	pullCmd := exec.Command("docker", "pull", fullImageName)
	if output, err := pullCmd.CombinedOutput(); err != nil {
		fmt.Printf("Error pulling image: %v\n%s\n", err, output)
		return
	}

	fmt.Printf("Successfully pulled image: %s\n", fullImageName)
}
