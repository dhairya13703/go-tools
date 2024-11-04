// cmd/cloudwatch/setup_alarm.go
package cloudwatch

import (
	"context"
	"fmt"

	"awesome-aws-cli/cmd/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/spf13/cobra"
)

func NewSetupAlarmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup-alarm",
		Short: "Set up a CloudWatch alarm",
		Run:   runSetupAlarm,
	}

	cmd.Flags().StringP("alarm-name", "a", "", "Alarm name")
	cmd.Flags().StringP("namespace", "n", "", "Metric namespace")
	cmd.Flags().StringP("metric-name", "m", "", "Metric name")
	cmd.Flags().Float64P("threshold", "t", 0, "Alarm threshold")
	cmd.MarkFlagRequired("alarm-name")
	cmd.MarkFlagRequired("namespace")
	cmd.MarkFlagRequired("metric-name")
	cmd.MarkFlagRequired("threshold")

	return cmd
}

func runSetupAlarm(cmd *cobra.Command, args []string) {
	alarmName, _ := cmd.Flags().GetString("alarm-name")
	namespace, _ := cmd.Flags().GetString("namespace")
	metricName, _ := cmd.Flags().GetString("metric-name")
	threshold, _ := cmd.Flags().GetFloat64("threshold")

	profile, _ := cmd.Flags().GetString("profile")
	region, _ := cmd.Flags().GetString("region")

	cfg, err := utils.LoadAWSConfig(profile, region)
	if err != nil {
		fmt.Printf("Error loading AWS config: %v\n", err)
		return
	}

	client := cloudwatch.NewFromConfig(cfg)

	input := &cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String(alarmName),
		ComparisonOperator: types.ComparisonOperatorGreaterThanThreshold,
		EvaluationPeriods:  aws.Int32(1),
		MetricName:         aws.String(metricName),
		Namespace:          aws.String(namespace),
		Period:             aws.Int32(60),
		Statistic:          types.StatisticAverage,
		Threshold:          aws.Float64(threshold),
		ActionsEnabled:     aws.Bool(true),
		AlarmDescription:   aws.String("Alarm created by CLI"),
		Unit:               types.StandardUnitNone,
	}

	_, err = client.PutMetricAlarm(context.TODO(), input)
	if err != nil {
		fmt.Printf("Error setting up alarm: %v\n", err)
		return
	}

	fmt.Printf("Alarm '%s' set up successfully\n", alarmName)
}
