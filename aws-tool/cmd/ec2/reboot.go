package ec2

import (
	"context"
	"fmt"
	"log"

	"awesome-aws-cli/cmd/utils"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
)

func NewRebootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reboot",
		Short: "Reboot an EC2 instance",
		Run:   runReboot,
	}

	cmd.Flags().StringP("instance-id", "i", "", "EC2 instance ID to reboot")

	return cmd
}

func runReboot(cmd *cobra.Command, args []string) {
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

	_, err = client.RebootInstances(context.TODO(), &ec2.RebootInstancesInput{
		InstanceIds: []string{instanceID},
	})

	if err != nil {
		log.Fatalf("Error rebooting instance %s: %v", instanceID, err)
	}

	fmt.Printf("Successfully initiated reboot for instance %s\n", instanceID)
}
