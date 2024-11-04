// cmd/cloudwatch/get_metric_data.go
package cloudwatch

import (
	"context"
	"fmt"
	"time"

	"awesome-aws-cli/cmd/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/spf13/cobra"
)

func NewGetMetricDataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-metric-data",
		Short: "Get CloudWatch metric data",
		Run:   runGetMetricData,
	}

	cmd.Flags().StringP("namespace", "n", "", "Metric namespace")
	cmd.Flags().StringP("metric-name", "m", "", "Metric name")
	cmd.Flags().StringP("start-time", "s", "", "Start time (format: 2006-01-02T15:04:05Z)")
	cmd.Flags().StringP("end-time", "e", "", "End time (format: 2006-01-02T15:04:05Z)")
	cmd.MarkFlagRequired("namespace")
	cmd.MarkFlagRequired("metric-name")

	return cmd
}

func runGetMetricData(cmd *cobra.Command, args []string) {
	namespace, _ := cmd.Flags().GetString("namespace")
	metricName, _ := cmd.Flags().GetString("metric-name")
	startTimeStr, _ := cmd.Flags().GetString("start-time")
	endTimeStr, _ := cmd.Flags().GetString("end-time")

	profile, _ := cmd.Flags().GetString("profile")
	region, _ := cmd.Flags().GetString("region")

	cfg, err := utils.LoadAWSConfig(profile, region)
	if err != nil {
		fmt.Printf("Error loading AWS config: %v\n", err)
		return
	}

	client := cloudwatch.NewFromConfig(cfg)

	var startTime, endTime time.Time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			fmt.Printf("Error parsing start time: %v\n", err)
			return
		}
	} else {
		startTime = time.Now().Add(-1 * time.Hour)
	}

	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			fmt.Printf("Error parsing end time: %v\n", err)
			return
		}
	} else {
		endTime = time.Now()
	}

	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []types.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &types.MetricStat{
					Metric: &types.Metric{
						Namespace:  aws.String(namespace),
						MetricName: aws.String(metricName),
					},
					Period: aws.Int32(60),
					Stat:   aws.String("Average"),
				},
			},
		},
		StartTime: aws.Time(startTime),
		EndTime:   aws.Time(endTime),
	}

	output, err := client.GetMetricData(context.TODO(), input)
	if err != nil {
		fmt.Printf("Error getting metric data: %v\n", err)
		return
	}

	for _, result := range output.MetricDataResults {
		fmt.Printf("Metric: %s\n", *result.Label)
		for i, timestamp := range result.Timestamps {
			fmt.Printf("  [%s] %f\n", timestamp.Format(time.RFC3339), result.Values[i])
		}
	}
}
