package cmd

import (
	"github.com/spf13/cobra"
)

func NewConfigCommand() *cobra.Command {
	ConfigCmd := &cobra.Command{
		Use: "config",
	}

	ConfigCmd.AddCommand(
		NewConfigGetCommand(),
		NewConfigSetCommand(),
	)
	return ConfigCmd
}
