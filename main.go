package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/cmd"
	"github.com/bluemir/wikinote/internal/buildinfo"

	// plugins
	_ "github.com/bluemir/wikinote/internal/plugins/__test__"
	_ "github.com/bluemir/wikinote/internal/plugins/discus"
	_ "github.com/bluemir/wikinote/internal/plugins/footer"
	_ "github.com/bluemir/wikinote/internal/plugins/recently-changes"
	_ "github.com/bluemir/wikinote/internal/plugins/utteranc.es"
)

var Version string
var AppName string

func main() {
	buildinfo.AppName = AppName
	buildinfo.Version = Version

	if err := cmd.Run(AppName, Version); err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
}
