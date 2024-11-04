package ec2

import (
	"context"
	"fmt"
	"log"

	"awesome-aws-cli/cmd/utils"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
)

func NewStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop an EC2 instance",
		Run:   runStop,
	}

	cmd.Flags().StringP("instance-id", "i", "", "EC2 instance ID to stop")

	return cmd
}

func runStop(cmd *cobra.Command, args []string) {
	instanceID, _ := cmd.Flags().GetString("instance-id")
	profile, _ := cmd.Flags().GetString("profile")
	region, _ := cmd.Flags().GetString("region")

	cfg, err := utils.LoadAWSConfig(profile, region)
	if err != nil {
		log.Fatalf("Error loading AWS config: %v", err)
	}

	client := ec2.NewFromConfig(cfg)

	if instanceID == "" {
		instanceID, err = listAndSelectInstance(client)
		if err != nil {
			log.Fatalf("Error selecting instance: %v", err)
		}
	}

	_, err = client.StopInstances(context.TODO(), &ec2.StopInstancesInput{
		InstanceIds: []string{instanceID},
	})

	if err != nil {
		log.Fatalf("Error stopping instance %s: %v", instanceID, err)
	}

	fmt.Printf("Successfully initiated stop for instance %s\n", instanceID)
}
