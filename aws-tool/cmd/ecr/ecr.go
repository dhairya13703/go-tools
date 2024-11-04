package ecr

import (
	"github.com/spf13/cobra"
)

func NewECRCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ecr",
		Short: "Manage ECR resources",
	}

	cmd.AddCommand(NewListRepositoriesCmd())
	cmd.AddCommand(NewGetLoginCmd())
	cmd.AddCommand(NewListImagesCmd())
	cmd.AddCommand(NewPushImageCmd())
	cmd.AddCommand(NewPullImageCmd())
	cmd.AddCommand(NewCreateRepositoryCmd())

	return cmd
}
