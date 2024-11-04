// cmd/cloudwatch/fetch_logs.go
package cloudwatch

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"awesome-aws-cli/cmd/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

func NewFetchLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fetch-logs",
		Short: "Fetch and display CloudWatch logs",
		Long: `Fetch and display CloudWatch logs. 
You can specify the start and end times in various formats:
- Just the date: "2024-10-22"
- Date and time: "2024-10-22 16:44"
- Date and time with seconds: "2024-10-22 16:44:28"
- ISO8601 format: "2024-10-22T16:44:28Z"
- RFC3339 format: "2024-10-22T16:44:28+00:00"
If no start time is specified, logs from the last 24 hours will be fetched by default.`,
		Run: runFetchLogs,
	}

	cmd.Flags().StringP("log-group", "g", "", "Log group name or number from list")
	cmd.Flags().StringP("log-stream", "l", "", "Log stream name")
	cmd.Flags().StringP("start-time", "s", "", "Start time (e.g., '2024-10-22' or '2024-10-22 16:44')")
	cmd.Flags().StringP("end-time", "e", "", "End time (e.g., '2024-10-22' or '2024-10-22 16:44')")
	cmd.Flags().IntP("hours", "H", 24, "Number of hours to look back if no start time is specified")

	return cmd
}

func runFetchLogs(cmd *cobra.Command, args []string) {
	logGroupInput, _ := cmd.Flags().GetString("log-group")
	logStreamInput, _ := cmd.Flags().GetString("log-stream")
	startTimeStr, _ := cmd.Flags().GetString("start-time")
	endTimeStr, _ := cmd.Flags().GetString("end-time")
	hours, _ := cmd.Flags().GetInt("hours")

	profile, _ := cmd.Flags().GetString("profile")
	region, _ := cmd.Flags().GetString("region")

	cfg, err := utils.LoadAWSConfig(profile, region)
	if err != nil {
		fmt.Printf("Error loading AWS config: %v\n", err)
		return
	}

	client := cloudwatchlogs.NewFromConfig(cfg)

	var logGroup string
	if logGroupInput == "" {
		// If no log group is specified, list log groups and prompt for selection
		listCmd := NewListLogGroupsCmd()
		listCmd.SetArgs([]string{})
		listCmd.Execute()

		logGroups, ok := listCmd.Context().Value("logGroups").([]string)
		if !ok || len(logGroups) == 0 {
			fmt.Println("No log groups available.")
			return
		}

		fmt.Print("Enter the number of the log group to fetch logs from: ")
		var selection string
		fmt.Scanln(&selection)

		index, err := strconv.Atoi(selection)
		if err != nil || index < 1 || index > len(logGroups) {
			fmt.Println("Invalid selection.")
			return
		}

		logGroup = logGroups[index-1]
	} else {
		logGroup = logGroupInput
	}

	var logStream string
	if logStreamInput == "" {
		// List log streams and prompt for selection
		streams, err := listLogStreams(client, logGroup)
		if err != nil {
			fmt.Printf("Error listing log streams: %v\n", err)
			return
		}

		if len(streams) == 0 {
			fmt.Println("No log streams found in the selected log group.")
			return
		}

		fmt.Println("Available log streams:")
		for i, stream := range streams {
			fmt.Printf("%d. %s\n", i+1, stream)
		}

		fmt.Print("Enter the number of the log stream to fetch logs from: ")
		var selection string
		fmt.Scanln(&selection)

		index, err := strconv.Atoi(selection)
		if err != nil || index < 1 || index > len(streams) {
			fmt.Println("Invalid selection.")
			return
		}

		logStream = streams[index-1]
	} else {
		logStream = logStreamInput
	}

	var startTime, endTime time.Time
	if startTimeStr != "" {
		startTime, err = parseFlexibleDateTime(startTimeStr)
		if err != nil {
			fmt.Printf("Error parsing start time: %v\n", err)
			return
		}
	} else {
		startTime = time.Now().Add(time.Duration(-hours) * time.Hour)
		fmt.Printf("Fetching logs from the last %d hours\n", hours)
	}

	if endTimeStr != "" {
		endTime, err = parseFlexibleDateTime(endTimeStr)
		if err != nil {
			fmt.Printf("Error parsing end time: %v\n", err)
			return
		}
	} else {
		endTime = time.Now()
	}

	// Check if start time is in the future
	if startTime.After(time.Now()) {
		fmt.Println("Warning: Start time is in the future. No logs will be available.")
		return
	}

	// Ensure end time is not in the future
	if endTime.After(time.Now()) {
		endTime = time.Now()
		fmt.Println("Warning: End time adjusted to current time.")
	}

	input := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroup),
		LogStreamName: aws.String(logStream),
		StartTime:     aws.Int64(startTime.UnixNano() / int64(time.Millisecond)),
		EndTime:       aws.Int64(endTime.UnixNano() / int64(time.Millisecond)),
		Limit:         aws.Int32(100), // Limit to 100 log events per request
	}

	paginator := cloudwatchlogs.NewGetLogEventsPaginator(client, input)

	fmt.Println("Fetching logs...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	eventCount := 0
	for paginator.HasMorePages() {
		select {
		case <-ctx.Done():
			fmt.Println("Operation timed out. Try narrowing your time range or increasing the timeout.")
			return
		default:
			output, err := paginator.NextPage(ctx)
			if err != nil {
				fmt.Printf("Error fetching logs: %v\n", err)
				return
			}

			for _, event := range output.Events {
				timestamp := time.Unix(*event.Timestamp/1000, 0)
				fmt.Printf("[%s] %s\n", timestamp.Format(time.RFC3339), *event.Message)
				eventCount++
			}

			if eventCount >= 1000 {
				fmt.Println("Reached 1000 log events. Stopping to prevent overwhelming output.")
				fmt.Println("To see more logs, please specify a narrower time range.")
				return
			}
		}
	}

	if eventCount == 0 {
		fmt.Println("No log events found in the specified time range.")
	} else {
		fmt.Printf("Retrieved %d log events.\n", eventCount)
	}
}

func listLogStreams(client *cloudwatchlogs.Client, logGroup string) ([]string, error) {
	input := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(logGroup),
	}

	paginator := cloudwatchlogs.NewDescribeLogStreamsPaginator(client, input)

	var streams []string
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		for _, stream := range output.LogStreams {
			streams = append(streams, *stream.LogStreamName)
		}
	}

	return streams, nil
}

func parseFlexibleDateTime(input string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		time.RFC3339,
	}

	input = strings.TrimSpace(input)

	for _, format := range formats {
		if t, err := time.Parse(format, input); err == nil {
			// If only date is provided, set time to beginning of the day (00:00:00)
			if len(input) == 10 {
				return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()), nil
			}
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date-time: %s", input)
}
