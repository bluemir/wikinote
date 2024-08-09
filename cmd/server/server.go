package server

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/cmd/config"
	"github.com/bluemir/wikinote/internal/backend"
	"github.com/bluemir/wikinote/internal/server"
)

func Register(cmd *kingpin.CmdClause, conf *config.Config) {
	cmd.Flag("bind", "bind").
		Default(":4000").
		StringVar(&conf.Server.Bind)
	cmd.Flag("front-page", "front page").
		Default("front-page.md").
		StringVar(&conf.Server.FrontPage)

	cmd.Action(func(*kingpin.ParseContext) error {
		logrus.Debugf("%#v", conf)

		ctx, stop := signal.NotifyContext(context.Background(),
			syscall.SIGTERM,
			syscall.SIGINT,
		)
		defer stop()

		b, err := backend.New(ctx, conf.Backend.Wikipath, conf.Backend.VolatileDatabase)
		if err != nil {
			logrus.Fatalf("%+v", err)
			return err
		}
		return server.Run(ctx, b, &conf.Server)
	})
}
