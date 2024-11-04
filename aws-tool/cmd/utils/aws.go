package utils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

func LoadAWSConfig(profile, region string) (aws.Config, error) {
	ctx := context.TODO()
	configOptions := []func(*config.LoadOptions) error{}

	if region != "" {
		configOptions = append(configOptions, config.WithRegion(region))
	}

	if profile != "" {
		configOptions = append(configOptions, config.WithSharedConfigProfile(profile))
	}

	return config.LoadDefaultConfig(ctx, configOptions...)
}

func GetECRRegistryURL(cfg aws.Config) (string, error) {
	client := ecr.NewFromConfig(cfg)
	input := &ecr.GetAuthorizationTokenInput{}
	resp, err := client.GetAuthorizationToken(context.TODO(), input)
	if err != nil {
		return "", err
	}

	if len(resp.AuthorizationData) == 0 {
		return "", fmt.Errorf("no authorization data received")
	}

	return *resp.AuthorizationData[0].ProxyEndpoint, nil
}
