package server

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
	"github.com/pkg/errors"
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
	cmd.Flag("https", "enable https mode").
		BoolVar(&conf.Server.EnableHttps)
	cmd.Flag("domain", "autoTLS domain").
		StringVar(&conf.Server.HttpsDomain)
	cmd.Flag("autotls-cache", "autotls cache").
		StringVar(&conf.Server.AutoTLSCache)

	cmd.Action(func(*kingpin.ParseContext) error {
		logrus.Debugf("%#v", conf)

		// validation
		if conf.Server.EnableHttps && conf.Server.HttpsDomain == "" {
			return errors.Errorf("https option need domain option")
		}
		if conf.Server.EnableHttps && conf.Server.Bind != ":4000" {
			return errors.Errorf("Bind option will be ignored")
		}

		ctx, stop := signal.NotifyContext(context.Background(),
			syscall.SIGTERM,
			syscall.SIGINT,
		)
		defer stop()

		b, err := backend.New(ctx, conf.Backend.Wikipath, conf.Backend.AdminUsers)
		if err != nil {
			logrus.Fatalf("%+v", err)
			return err
		}
		return server.Run(ctx, b, &conf.Server)
	})
}
