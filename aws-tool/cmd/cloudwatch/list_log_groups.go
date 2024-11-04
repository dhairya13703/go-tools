package cloudwatch

import (
	"context"
	"fmt"
	"os"
	"time"

	"awesome-aws-cli/cmd/utils"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewListLogGroupsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-log-groups",
		Short: "List CloudWatch log groups",
		Run:   runListLogGroups,
	}

	return cmd
}

func runListLogGroups(cmd *cobra.Command, args []string) {
	profile, _ := cmd.Flags().GetString("profile")
	region, _ := cmd.Flags().GetString("region")

	cfg, err := utils.LoadAWSConfig(profile, region)
	if err != nil {
		fmt.Printf("Error loading AWS config: %v\n", err)
		return
	}

	client := cloudwatchlogs.NewFromConfig(cfg)

	input := &cloudwatchlogs.DescribeLogGroupsInput{}
	paginator := cloudwatchlogs.NewDescribeLogGroupsPaginator(client, input)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "Log Group Name", "Creation Time", "Retention (days)"})

	var logGroups []string
	index := 1

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			fmt.Printf("Error listing log groups: %v\n", err)
			return
		}

		for _, logGroup := range output.LogGroups {
			logGroups = append(logGroups, *logGroup.LogGroupName)
			retention := "Never Expire"
			if logGroup.RetentionInDays != nil {
				retention = fmt.Sprintf("%d", *logGroup.RetentionInDays)
			}
			table.Append([]string{
				fmt.Sprintf("%d", index),
				*logGroup.LogGroupName,
				time.Unix(*logGroup.CreationTime/1000, 0).String(),
				retention,
			})
			index++
		}
	}

	table.Render()

	// Store log groups for later use
	cmd.SetContext(context.WithValue(cmd.Context(), "logGroups", logGroups))
}
