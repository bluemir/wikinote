package cmd

import (
	"github.com/spf13/cobra"
)

func NewUserRemoveCommand() *cobra.Command {
	var UserRemoveCmd = &cobra.Command{
		Use: "rm",
	}
	return UserRemoveCmd
}
