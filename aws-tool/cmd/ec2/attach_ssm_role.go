package ec2

import (
	"context"
	"fmt"
	"log"
	"time"

	"awesome-aws-cli/cmd/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/spf13/cobra"
)

const (
	ssmRoleName            = "SSMRoleForEC2"
	ssmInstanceProfileName = "SSMInstanceProfileForEC2"
	ssmPolicyARN           = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
	assumeRolePolicyDoc    = `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"Service": "ec2.amazonaws.com"
				},
				"Action": "sts:AssumeRole"
			}
		]
	}`
)

func NewAttachSSMRoleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "attach-ssm-role -i [instance-id]",
		Short: "Attach SSM role to an EC2 instance",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			instanceID, _ := cmd.Flags().GetString("instance-id")
			if instanceID == "" {
				log.Fatal("Instance ID is required")
			}
			profile, _ := cmd.Flags().GetString("profile")
			region, _ := cmd.Flags().GetString("region")

			cfg, err := utils.LoadAWSConfig(profile, region)
			if err != nil {
				log.Fatalf("Error loading AWS config: %v", err)
			}

			iamClient := iam.NewFromConfig(cfg)
			ec2Client := ec2.NewFromConfig(cfg)

			// Check if the role exists, create if it doesn't
			role, err := getOrCreateSSMRole(iamClient)
			if err != nil {
				log.Fatalf("Error getting or creating SSM role: %v", err)
			}

			// Create instance profile if it doesn't exist
			instanceProfile, err := getOrCreateInstanceProfile(iamClient, *role.Role.RoleName)
			if err != nil {
				log.Fatalf("Error getting or creating instance profile: %v", err)
			}

			// Wait for the instance profile to be ready
			time.Sleep(10 * time.Second)

			// Attach role to instance
			err = attachRoleToInstance(ec2Client, instanceID, *instanceProfile.InstanceProfile.InstanceProfileName)
			if err != nil {
				log.Fatalf("Error attaching role to instance: %v", err)
			}

			fmt.Printf("Successfully attached SSM role to instance %s\n", instanceID)
		},
	}
}

func init() {
	attachSSMRoleCmd := NewAttachSSMRoleCmd()
	attachSSMRoleCmd.Flags().StringP("instance-id", "i", "", "EC2 instance ID to attach the SSM role to")
	attachSSMRoleCmd.MarkFlagRequired("instance-id")
}

func getOrCreateSSMRole(client *iam.Client) (*iam.GetRoleOutput, error) {
	role, err := client.GetRole(context.TODO(), &iam.GetRoleInput{
		RoleName: aws.String(ssmRoleName),
	})

	if err == nil {
		return role, nil
	}

	// If the role doesn't exist, create it
	createRoleOutput, err := client.CreateRole(context.TODO(), &iam.CreateRoleInput{
		RoleName:                 aws.String(ssmRoleName),
		AssumeRolePolicyDocument: aws.String(assumeRolePolicyDoc),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	_, err = client.AttachRolePolicy(context.TODO(), &iam.AttachRolePolicyInput{
		RoleName:  aws.String(ssmRoleName),
		PolicyArn: aws.String(ssmPolicyARN),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to attach policy to role: %w", err)
	}

	return &iam.GetRoleOutput{Role: createRoleOutput.Role}, nil
}

func getOrCreateInstanceProfile(client *iam.Client, roleName string) (*iam.GetInstanceProfileOutput, error) {
	instanceProfile, err := client.GetInstanceProfile(context.TODO(), &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(ssmInstanceProfileName),
	})

	if err == nil {
		return instanceProfile, nil
	}

	// If the instance profile doesn't exist, create it
	createProfileOutput, err := client.CreateInstanceProfile(context.TODO(), &iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String(ssmInstanceProfileName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create instance profile: %w", err)
	}

	_, err = client.AddRoleToInstanceProfile(context.TODO(), &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: aws.String(ssmInstanceProfileName),
		RoleName:            aws.String(roleName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add role to instance profile: %w", err)
	}

	return &iam.GetInstanceProfileOutput{InstanceProfile: createProfileOutput.InstanceProfile}, nil
}

func attachRoleToInstance(client *ec2.Client, instanceID, profileName string) error {
	_, err := client.AssociateIamInstanceProfile(context.TODO(), &ec2.AssociateIamInstanceProfileInput{
		IamInstanceProfile: &types.IamInstanceProfileSpecification{
			Name: aws.String(profileName),
		},
		InstanceId: aws.String(instanceID),
	})
	if err != nil {
		return fmt.Errorf("failed to associate IAM instance profile: %w", err)
	}
	return nil
}
