package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bluemir/wikinote/pkg/backend"
	"github.com/bluemir/wikinote/pkg/server"

	// plugins
	_ "github.com/bluemir/wikinote/pkg/plugins/discus"
	_ "github.com/bluemir/wikinote/pkg/plugins/recently-changes"
)

var VERSION string

type Config struct {
	Backend backend.Config
	Server  server.Config
	Debug   bool
}

func main() {
	// log
	if level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL")); err != nil {
		logrus.Warn("unknown log level. using default level(info)")
	} else {
		logrus.SetLevel(level)
	}

	conf := &Config{}

	cli := kingpin.New("wikinote", "main code")

	cli.Flag("debug", "enable debug mode").BoolVar(&conf.Debug)
	cli.Flag("wiki-path", "wiki data path").
		Short('w').
		Default(os.ExpandEnv("$HOME/wiki")).
		PlaceHolder("$HOME/wiki").
		StringVar(&conf.Backend.Wikipath)
	cli.Flag("config", "config file").
		Short('c').
		Default(os.ExpandEnv("$HOME/wiki/.app/config.yaml")).
		PlaceHolder("$HOME/wiki/.app/config.yaml").
		StringVar(&conf.Backend.ConfigFile)
	cli.Flag("role", "role file").
		Short('r').
		Default(os.ExpandEnv("$HOME/wiki/.app/roles.yaml")).
		PlaceHolder("$HOME/wiki/.app/roles.yaml").
		StringVar(&conf.Backend.RoleFile)
	cli.Flag("root-user", "root user").
		StringVar(&conf.Backend.RootUser)

	cli.Version(VERSION)

	serve := cli.Command("serve", "run server")
	{
		serve.Flag("bind", "bind address").Default(":4000").StringVar(&conf.Server.Bind)
		serve.Flag("tls-domain", "tls domain").StringsVar(&conf.Server.TLSDomains)
	}

	switch kingpin.MustParse(cli.Parse(os.Args[1:])) {
	case serve.FullCommand():
		logrus.Debugf("%#v", conf)
		b, err := backend.New(&conf.Backend)
		if err != nil {
			logrus.Fatal(err)
		}
		if err := server.Run(b, &conf.Server); err != nil {
			logrus.Fatal(err)
		}
	}

}
