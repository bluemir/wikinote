package cmd

import (
	"github.com/spf13/cobra"
)

// userCmd represents the user command

func NewUserCommand() *cobra.Command {
	var UserCmd = &cobra.Command{
		Use:   "user",
		Short: "manage user",
	}
	UserCmd.AddCommand(
		NewUserListCommand(),
		NewUserAddCommand(),
		NewUserRemoveCommand(),
		NewUserAssignCommand(),
		NewUserPasswordCommand(),
	)
	return UserCmd
}
