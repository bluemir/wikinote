package cmd

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/bluemir/wikinote/backend"
	"github.com/bluemir/wikinote/server"
)

func NewServeCommand() *cobra.Command {
	var ServeConfig = struct {
		port uint
		bind string
		opts []string
	}{}

	var ServeCmd = &cobra.Command{
		Use:   "serve",
		Short: "A brief description of your application",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := &server.Config{}
			if cmd.Flags().Changed("port") {
				conf.Address = fmt.Sprintf(":%d", ServeConfig.port)
			} else {
				conf.Address = ServeConfig.bind
			}
			b, err := backend.New(backendOpts)
			if err != nil {
				return err
			}

			for k, v := range parseKeyValue(ServeConfig.opts) {
				logrus.Debugf("key:'%s' value:'%s'", k, v)
				b.Config().Set(k, v)
			}
			return server.Run(b, conf)
		},
	}

	ServeCmd.Flags().UintVarP(&ServeConfig.port, "port", "p", 0, "listen port")
	ServeCmd.Flags().StringVarP(&ServeConfig.bind, "bind", "b", "localhost:4000", "bind address")
	ServeCmd.Flags().StringArrayVarP(&ServeConfig.opts, "opt", "o", []string{}, "live options")

	return ServeCmd
}
func parseKeyValue(pairs []string) map[string]string {
	result := map[string]string{}
	for _, kv := range pairs {
		s := strings.SplitN(kv, "=", 2)
		if len(s) < 2 {
			s = append(s, "")
		}

		result[s[0]] = s[1]
	}
	return result
}
