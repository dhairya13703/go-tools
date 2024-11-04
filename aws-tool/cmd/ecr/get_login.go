package ecr

import (
	"context"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"

	"awesome-aws-cli/cmd/utils"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/spf13/cobra"
)

func NewGetLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-login",
		Short: "Get docker login command",
		Run:   runGetLogin,
	}

	return cmd
}

func runGetLogin(cmd *cobra.Command, args []string) {
	profile, _ := cmd.Flags().GetString("profile")
	region, _ := cmd.Flags().GetString("region")

	cfg, err := utils.LoadAWSConfig(profile, region)
	if err != nil {
		fmt.Printf("Error loading AWS config: %v\n", err)
		return
	}

	client := ecr.NewFromConfig(cfg)

	input := &ecr.GetAuthorizationTokenInput{}
	resp, err := client.GetAuthorizationToken(context.TODO(), input)
	if err != nil {
		fmt.Printf("Error getting authorization token: %v\n", err)
		return
	}

	if len(resp.AuthorizationData) > 0 {
		authToken := *resp.AuthorizationData[0].AuthorizationToken
		decodedToken, err := base64.StdEncoding.DecodeString(authToken)
		if err != nil {
			fmt.Printf("Error decoding authorization token: %v\n", err)
			return
		}

		parts := strings.SplitN(string(decodedToken), ":", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid authorization token format")
			return
		}

		// username := parts[0]
		password := parts[1]

		proxyEndpoint := *resp.AuthorizationData[0].ProxyEndpoint

		// Execute the login command
		loginCmd := exec.Command("docker", "login",
			"--username", "AWS",
			"--password-stdin",
			proxyEndpoint)

		loginCmd.Stdin = strings.NewReader(password)
		output, err := loginCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error logging in: %v\n", err)
			return
		}

		fmt.Println(string(output))
	} else {
		fmt.Println("No authorization data received")
	}
}
