/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package ec2

import (
	"github.com/spf13/cobra"
)

func NewEC2Cmd() *cobra.Command {
	ec2Cmd := &cobra.Command{
		Use:   "ec2",
		Short: "Manage EC2 instances",
		Long:  `Perform various operations on EC2 instances such as list, start, stop, reboot, and attach SSM role.`,
	}

	// Add EC2 subcommands
	ec2Cmd.AddCommand(NewListCmd())
	ec2Cmd.AddCommand(NewAttachSSMRoleCmd())
	ec2Cmd.AddCommand(NewRebootCmd())
	ec2Cmd.AddCommand(NewStartCmd())
	ec2Cmd.AddCommand(NewStopCmd())
	ec2Cmd.AddCommand(NewListEBSVolumesCmd())

	return ec2Cmd
}
