package cmd

import (
	"github.com/spf13/cobra"
)

func NewUserAssignCommand() *cobra.Command {
	var UserAssignCmd = &cobra.Command{
		Use: "user COMMAND",
	}
	return UserAssignCmd
}
