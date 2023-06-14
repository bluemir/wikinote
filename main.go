package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/cmd"
	"github.com/bluemir/wikinote/internal/buildinfo"

	// plugins
	_ "github.com/bluemir/wikinote/internal/plugins/__test__"
	_ "github.com/bluemir/wikinote/internal/plugins/footer"
	_ "github.com/bluemir/wikinote/internal/plugins/giscus"
	_ "github.com/bluemir/wikinote/internal/plugins/last-modified"
	_ "github.com/bluemir/wikinote/internal/plugins/recently-changes"
)

var Version string
var AppName string
var BuildTime string

func main() {
	buildinfo.AppName = AppName
	buildinfo.Version = Version
	buildinfo.AppName = BuildTime

	if err := cmd.Run(AppName, Version); err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
}
