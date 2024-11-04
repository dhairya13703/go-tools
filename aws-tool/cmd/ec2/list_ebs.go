package ec2

import (
	"awesome-aws-cli/cmd/utils"
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewListEBSVolumesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-volumes [instance-id]",
		Short: "List EBS volumes",
		Run:   runListEBSVolumes,
	}

	return cmd
}

func runListEBSVolumes(cmd *cobra.Command, args []string) {
	fmt.Println("Listing EBS Volumes...")

	profile, _ := cmd.Flags().GetString("profile")
	region, _ := cmd.Flags().GetString("region")

	cfg, err := utils.LoadAWSConfig(profile, region)
	if err != nil {
		fmt.Printf("Error loading AWS config: %v\n", err)
		return
	}

	client := ec2.NewFromConfig(cfg)

	// Fix this input error

	var input *ec2.DescribeVolumesInput
	if len(args) > 0 {
		instanceID := args[0]
		input = &ec2.DescribeVolumesInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("attachment.instance-id"),
					Values: []string{instanceID},
				},
			},
		}
	} else {
		input = &ec2.DescribeVolumesInput{}
	}

	resp, err := client.DescribeVolumes(context.TODO(), input)
	if err != nil {
		fmt.Printf("Unable to describe EBS: %v\n", err)
		return
	}

	if len(resp.Volumes) == 0 {
		fmt.Println("No EBS volumes found.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Volume ID", "Size (GB)", "Volume Type", "State", "Attached Instance ID", "Device Name"})

	for _, volume := range resp.Volumes {
		var instanceID, deviceName string

		if len(volume.Attachments) > 0 {
			instanceID = *volume.Attachments[0].InstanceId
			deviceName = *volume.Attachments[0].Device
		}

		table.Append([]string{
			*volume.VolumeId,
			strconv.FormatInt(int64(*volume.Size), 10),
			string(volume.VolumeType),
			string(volume.State),
			instanceID,
			deviceName,
		})
	}
	table.Render()

}
