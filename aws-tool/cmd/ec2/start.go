package ec2

import (
	"context"
	"fmt"
	"log"

	"awesome-aws-cli/cmd/utils"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
)

func NewStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start an EC2 instance",
		Run:   runStart,
	}

	cmd.Flags().StringP("instance-id", "i", "", "EC2 instance ID to start")

	return cmd
}

func runStart(cmd *cobra.Command, args []string) {
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

	_, err = client.StartInstances(context.TODO(), &ec2.StartInstancesInput{
		InstanceIds: []string{instanceID},
	})

	if err != nil {
		log.Fatalf("Error starting instance %s: %v", instanceID, err)
	}

	fmt.Printf("Successfully initiated start for instance %s\n", instanceID)
}

func newStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start [instance-id]",
		Short: "Start an EC2 instance",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			profile, _ := cmd.Flags().GetString("profile")
			region, _ := cmd.Flags().GetString("region")
			cfg, err := utils.LoadAWSConfig(profile, region)
			if err != nil {
				fmt.Printf("Error loading AWS config: %v\n", err)
				return
			}

			client := ec2.NewFromConfig(cfg)

			_, err = client.StartInstances(context.TODO(), &ec2.StartInstancesInput{
				InstanceIds: []string{args[0]},
			})
			if err != nil {
				fmt.Printf("Unable to start instance: %v\n", err)
				return
			}

			fmt.Printf("Successfully started instance %s\n", args[0])
		},
	}
}
