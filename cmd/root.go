package cmd

import (
	"errors"
	"os"

	docopt "github.com/docopt/docopt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Execute(version string) error {
	// log
	if level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL")); err != nil {
		logrus.Warn("unknown log level. using default level(info)")
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(level)
	}

	//docopt.Parse(doc, argv, help, version, optionsFirst)
	args, err := docopt.Parse(usage, os.Args[1:], true, version, true)
	if err != nil {
		logrus.Panicf("error on parse usage %s", err)
		return err
	}

	logrus.Debug(args)
	if args["--debug"].(bool) {
		logrus.Info("Turn on debug mode")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	switch args["<command>"] {
	case "serve":
		if err := doServe(args["<args>"].([]string), version); err != nil {
			return err
		}
	case "user", "config":
		// proxy
		return errors.New("Not Implements")
	default:
		return errors.New("Not Implements")
	}
	return nil
}
