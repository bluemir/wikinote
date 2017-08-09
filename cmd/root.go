package cmd

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/bluemir/wikinote/backend"
)

var backendOpts = &backend.Options{}

func Execute() {
	var RootConfig = struct {
		debug bool
	}{}
	var RootCmd = &cobra.Command{
		Use:   "wikinote",
		Short: "A simple markdown renderer",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if RootConfig.debug {
				logrus.Info("Turn on debug mode")
				logrus.SetLevel(logrus.DebugLevel)
			} else {
				gin.SetMode(gin.ReleaseMode)
			}
		},
	}
	RootCmd.PersistentFlags().BoolVarP(&RootConfig.debug, "debug", "D", false, "debug mode")
	RootCmd.PersistentFlags().StringVarP(&backendOpts.Wikipath, "wiki-path", "w", "$HOME/wiki", "wikipath")
	RootCmd.PersistentFlags().StringVarP(&backendOpts.ConfigFile, "config", "c", "$HOME/wiki/.app/config.yaml", "config file")

	RootCmd.AddCommand(
		NewUserCommand(),
		NewConfigCommand(),
		NewServeCommand(),
	)

	// make backend?
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
