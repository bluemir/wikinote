package cmd

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/pkgs/backend"
	"github.com/bluemir/wikinote/pkgs/server"
)

func doServe(argv []string, version string) error {
	argv = append([]string{"serve"}, argv...)
	logrus.Debug(argv)
	args, err := docopt.Parse(serveUsage, argv, true, version, false)
	if err != nil {
		logrus.Errorf("%+q", err)
		return err
	}
	logrus.Debug(args)

	opts := &backend.Options{
		ConfigFile: args["--config"].(string),
		Wikipath:   args["--wiki-path"].(string),
		Version:    version,
	}
	b, err := backend.New(opts)
	if err != nil {
		return err
	}

	conf := &server.Config{
		Address: args["--bind"].(string),
		//Domain:  []string{args["--tls-domain"].(string)},
	}

	if domain, ok := args["--tls-domain"].(string); ok {
		conf.Domain = []string{domain}
	}

	return server.Run(b, conf)
}
