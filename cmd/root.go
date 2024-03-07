package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bluemir/wikinote/cmd/config"
	"github.com/bluemir/wikinote/internal/buildinfo"

	serverCmd "github.com/bluemir/wikinote/cmd/server"
)

const (
	describe        = ``
	defaultLogLevel = logrus.InfoLevel
)

func Run() error {
	conf := config.NewConfig()

	app := kingpin.New(buildinfo.AppName, describe)
	app.Version(buildinfo.Version)

	app.Flag("verbose", "Log level").
		Short('v').
		CounterVar(&conf.LogLevel)
	app.Flag("log-format", "Log format").
		StringVar(&conf.LogFormat)
	app.PreAction(setupLogger(conf))

	// app flags
	app.Flag("wiki-path", "wiki data path").
		Short('w').
		Default(os.ExpandEnv("$HOME/wiki")).
		PlaceHolder("$HOME/wiki").
		StringVar(&conf.Backend.Wikipath)
	app.Flag("admin-user", "admin user").
		StringMapVar(&conf.Backend.AdminUsers)

	// server flags

	serverCmd.Register(app.Command("server", "server"), conf)
	//clientCmd.Register(app.Command("client", "client"))

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	logrus.Debug("Shutdown")

	return nil
}
func setupLogger(conf *config.Config) func(*kingpin.ParseContext) error {
	return func(*kingpin.ParseContext) error {
		level := logrus.Level(conf.LogLevel) + defaultLogLevel
		logrus.SetOutput(os.Stderr)
		logrus.SetLevel(level)
		logrus.SetReportCaller(true)
		logrus.Infof("logrus level: %s", level)

		switch conf.LogFormat {
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
	}
}
