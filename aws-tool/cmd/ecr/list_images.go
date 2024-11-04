package ecr

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"awesome-aws-cli/cmd/utils"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewListImagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-images [repository-name]",
		Short: "List images in an ECR repository",
		Args:  cobra.ExactArgs(1),
		Run:   runListImages,
	}

	return cmd
}

func runListImages(cmd *cobra.Command, args []string) {
	repositoryName := args[0]
	profile, _ := cmd.Flags().GetString("profile")
	region, _ := cmd.Flags().GetString("region")

	cfg, err := utils.LoadAWSConfig(profile, region)
	if err != nil {
		fmt.Printf("Error loading AWS config: %v\n", err)
		return
	}

	client := ecr.NewFromConfig(cfg)

	input := &ecr.DescribeImagesInput{
		RepositoryName: &repositoryName,
	}

	paginator := ecr.NewDescribeImagesPaginator(client, input)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Tag", "Digest", "Size (MB)", "Pushed At", "Image URI"})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			fmt.Printf("Error listing images: %v\n", err)
			return
		}

		for _, imageDetail := range output.ImageDetails {
			tag := "N/A"
			if len(imageDetail.ImageTags) > 0 {
				tag = imageDetail.ImageTags[0]
			}

			size := float64(*imageDetail.ImageSizeInBytes) / 1024 / 1024 // Convert to MB
			sizeStr := strconv.FormatFloat(size, 'f', 2, 64)

			pushedAt := imageDetail.ImagePushedAt.Format(time.RFC3339)

			registryId := *imageDetail.RegistryId
			imageUri := fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/%s:%s", registryId, region, repositoryName, tag)

			table.Append([]string{tag, *imageDetail.ImageDigest, sizeStr, pushedAt, imageUri})
		}
	}

	table.Render()
}
