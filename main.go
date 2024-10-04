package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/cmd"
	// plugins
	_ "github.com/bluemir/wikinote/internal/plugins/_sample"
	//_ "github.com/bluemir/wikinote/internal/plugins/footer"
	//_ "github.com/bluemir/wikinote/internal/plugins/giscus"
	//_ "github.com/bluemir/wikinote/internal/plugins/last-modified"
	//_ "github.com/bluemir/wikinote/internal/plugins/recently-changes"
)

func main() {

	if err := cmd.Run(); err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
}
