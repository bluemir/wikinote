package main

import (
	"github.com/bluemir/wikinote/cmd"
	"github.com/sirupsen/logrus"
)

var Version string

func main() {

	if err := cmd.Execute(Version); err != nil {
		logrus.Error(err)
	}
}
