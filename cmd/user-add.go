package cmd

import (
	"github.com/spf13/cobra"
)

func NewUserAddCommand() *cobra.Command {
	var UserAddConfig = struct {
		email    string
		password string
	}{}
	var UserAddCmd = &cobra.Command{
		Use: "add USERNAME",
		RunE: func(cmd *cobra.Command, args []string) error {
			//b := backend.New(backendOpts)
			return nil
		},
	}

	UserAddCmd.Flags().StringVar(&UserAddConfig.email, "email", "USERNAME@wikinote", "user email")

	return UserAddCmd
}
