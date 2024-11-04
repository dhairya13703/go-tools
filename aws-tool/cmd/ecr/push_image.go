package ecr

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"awesome-aws-cli/cmd/utils"

	"github.com/spf13/cobra"
)

func NewPushImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push [repository-name] [local-image-name:tag]",
		Short: "Push an image to ECR",
		Args:  cobra.ExactArgs(2),
		Run:   runPushImage,
	}

	return cmd
}

func runPushImage(cmd *cobra.Command, args []string) {
	repositoryName := args[0]
	localImageName := args[1]
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

	// Remove the "https://" prefix from the registry URL
	registryURL = strings.TrimPrefix(registryURL, "https://")

	// Construct the full image name for ECR
	fullImageName := fmt.Sprintf("%s/%s:%s", registryURL, repositoryName, strings.Split(localImageName, ":")[1])

	fmt.Printf("Tagging image %s as %s\n", localImageName, fullImageName)

	// Tag the image
	tagCmd := exec.Command("docker", "tag", localImageName, fullImageName)
	if output, err := tagCmd.CombinedOutput(); err != nil {
		fmt.Printf("Error tagging image: %v\n%s\n", err, output)
		return
	}

	fmt.Printf("Pushing image to ECR: %s\n", fullImageName)

	// Push the image
	pushCmd := exec.Command("docker", "push", fullImageName)
	stdout, err := pushCmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating StdoutPipe: %v\n", err)
		return
	}

	if err := pushCmd.Start(); err != nil {
		fmt.Printf("Error starting push command: %v\n", err)
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var result map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &result); err != nil {
			fmt.Println(scanner.Text())
		} else {
			if status, ok := result["status"].(string); ok {
				fmt.Print(status)
				if progress, ok := result["progress"].(string); ok {
					fmt.Printf(": %s", progress)
				}
				fmt.Println()
			}
		}
	}

	if err := pushCmd.Wait(); err != nil {
		fmt.Printf("Error pushing image: %v\n", err)
		return
	}

	fmt.Printf("Successfully pushed image: %s\n", fullImageName)
}
