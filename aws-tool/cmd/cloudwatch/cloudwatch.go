package cloudwatch

import (
	"github.com/spf13/cobra"
)

func NewCloudWatchCmd() *cobra.Command {
	cwCmd := &cobra.Command{
		Use:   "cloudwatch",
		Short: "Manage CloudWatch logs, metrics, and alarms",
	}
	
	cwCmd.AddCommand(NewListLogGroupsCmd())
	cwCmd.AddCommand(NewFetchLogsCmd())
	cwCmd.AddCommand(NewGetMetricDataCmd())
	cwCmd.AddCommand(NewSetupAlarmCmd())

	return cwCmd
}
