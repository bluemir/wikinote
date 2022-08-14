package server

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

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
		b, err := backend.New(conf.Backend.Wikipath, conf.Backend.AdminUsers)
		if err != nil {
			logrus.Fatalf("%+v", err)
			return err
		}
		return server.Run(b, &conf.Server)
	})
}

// - http
//     - bind
// - https
//     - domain
