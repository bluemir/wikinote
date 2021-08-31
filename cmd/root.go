package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bluemir/wikinote/internal/backend"
	"github.com/bluemir/wikinote/internal/server"
)

const (
	describe        = ``
	defaultLogLevel = logrus.InfoLevel
)

func Run(AppName string, Version string) error {
	conf := struct {
		Backend   backend.Config
		Server    server.Config
		logLevel  int
		logFormat string
	}{
		Backend: backend.InitConfig(),
	}

	app := kingpin.New(AppName, describe)
	app.Version(Version)

	app.Flag("verbose", "Log level").
		Short('v').
		CounterVar(&conf.logLevel)
	app.Flag("log-format", "Log format").
		StringVar(&conf.logFormat)
	app.PreAction(func(*kingpin.ParseContext) error {
		level := logrus.Level(conf.logLevel) + defaultLogLevel
		logrus.SetOutput(os.Stderr)
		logrus.SetLevel(level)
		logrus.SetReportCaller(true)
		logrus.Infof("logrus level: %s", level)

		switch conf.logFormat {
		case "text-color":
			logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
		case "text":
			logrus.SetFormatter(&logrus.TextFormatter{})
		case "json":
			logrus.SetFormatter(&logrus.JSONFormatter{})
		case "":
			// do nothing. it means smart.
		default:
			return errors.Errorf("unknown log format")
		}

		return nil
	})

	// app flags
	app.Flag("wiki-path", "wiki data path").
		Short('w').
		Default(os.ExpandEnv("$HOME/wiki")).
		PlaceHolder("$HOME/wiki").
		StringVar(&conf.Backend.Wikipath)
	app.Flag("config", "config file").
		Short('c').
		Default(os.ExpandEnv("$HOME/wiki/.app/config.yaml")).
		PlaceHolder("$HOME/wiki/.app/config.yaml").
		StringVar(&conf.Backend.ConfigFile)
	app.Flag("admin-user", "admin user").
		StringMapVar(&conf.Backend.AdminUsers)

	// server flags
	serverCmd := app.Command("server", "server")
	{
		serverCmd.Flag("bind", "bind").
			Default(":4000").
			StringVar(&conf.Server.Bind)

		serverCmd.Flag("tls-domain", "tls domain").
			StringsVar(&conf.Server.TLSDomains)
	}

	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	switch cmd {
	case serverCmd.FullCommand():
		logrus.Debugf("%#v", conf)
		b, err := backend.New(&conf.Backend)
		if err != nil {
			logrus.Fatal(err)
			return err
		}
		return server.Run(b, &conf.Server)
	default:
		return errors.New("not implements command")
	}
}
