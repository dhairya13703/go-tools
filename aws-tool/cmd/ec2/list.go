package ec2

import (
	"context"
	"fmt"
	"os"

	"awesome-aws-cli/cmd/utils"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// NewListCmd creates a new command for listing EC2 instances
func NewListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List EC2 instances",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Listing EC2 instances...")

			profile, _ := cmd.Flags().GetString("profile")
			region, _ := cmd.Flags().GetString("region")

			cfg, err := utils.LoadAWSConfig(profile, region)
			if err != nil {
				fmt.Printf("Error loading AWS config: %v\n", err)
				return
			}

			client := ec2.NewFromConfig(cfg)

			resp, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
			if err != nil {
				fmt.Printf("Unable to describe instances: %v\n", err)
				return
			}

			if len(resp.Reservations) == 0 {
				fmt.Println("No EC2 instances found.")
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Instance ID", "Instance Type", "State", "Public IP", "Private IP"})

			for _, reservation := range resp.Reservations {
				for _, instance := range reservation.Instances {
					name := getInstanceName(instance.Tags)
					publicIP := ""
					if instance.PublicIpAddress != nil {
						publicIP = *instance.PublicIpAddress
					}
					privateIP := ""
					if instance.PrivateIpAddress != nil {
						privateIP = *instance.PrivateIpAddress
					}

					table.Append([]string{
						name,
						*instance.InstanceId,
						string(instance.InstanceType),
						string(instance.State.Name),
						publicIP,
						privateIP,
					})
				}
			}

			table.Render()
		},
	}
}
