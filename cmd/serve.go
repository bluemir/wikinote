package cmd

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/backend"
	"github.com/bluemir/wikinote/server"
)

func doServe(argv []string, version string) error {

	args, err := docopt.Parse(serveUsage, append([]string{"serve"}, argv...), true, version, true)
	if err != nil {
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

	return server.Run(b, conf)
}
