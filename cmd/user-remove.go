package cmd

import (
	"fmt"

	"github.com/bluemir/wikinote/backend"
	"github.com/spf13/cobra"
)

func NewUserRemoveCommand() *cobra.Command {
	var UserRemoveCmd = &cobra.Command{
		Use:     "remove USER",
		Aliases: []string{"rm", "del", "delete"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Invalid Arguments")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]

			b, err := backend.New(backendOpts)
			if err != nil {
				return err
			}

			return b.User().Delete(username)
		},
	}
	return UserRemoveCmd
}
