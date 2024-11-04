package ec2

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/olekukonko/tablewriter"
)

func listAndSelectInstance(client *ec2.Client) (string, error) {
	resp, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return "", fmt.Errorf("failed to describe instances: %w", err)
	}

	var instances []types.Instance
	for _, reservation := range resp.Reservations {
		instances = append(instances, reservation.Instances...)
	}

	if len(instances) == 0 {
		return "", fmt.Errorf("no instances found")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Number", "Instance ID", "Name", "State", "Public IP", "Private IP"})

	for i, instance := range instances {
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
			strconv.Itoa(i + 1),
			*instance.InstanceId,
			name,
			string(instance.State.Name),
			publicIP,
			privateIP,
		})
	}

	table.Render()

	var selection string
	fmt.Print("Enter the number of the instance to select (or 'q' to quit): ")
	fmt.Scanln(&selection)

	if selection == "q" {
		return "", fmt.Errorf("operation cancelled by user")
	}

	index, err := strconv.Atoi(selection)
	if err != nil || index < 1 || index > len(instances) {
		return "", fmt.Errorf("invalid selection")
	}

	return *instances[index-1].InstanceId, nil
}

func getInstanceName(tags []types.Tag) string {
	for _, tag := range tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return ""
}
