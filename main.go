package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/cmd"

	// plugins
	_ "github.com/bluemir/wikinote/pkg/plugins/__test__"
	_ "github.com/bluemir/wikinote/pkg/plugins/discus"
	_ "github.com/bluemir/wikinote/pkg/plugins/recently-changes"
	_ "github.com/bluemir/wikinote/pkg/plugins/utteranc.es"
)

var Version string
var AppName string

func main() {
	if err := cmd.Run(AppName, Version); err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
}
