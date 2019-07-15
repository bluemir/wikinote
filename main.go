package main

import (
	"github.com/bluemir/wikinote/cmd"
	"github.com/sirupsen/logrus"
)

var Version string

func main() {

	err := cmd.Execute(Version)
	if err != nil {
		logrus.Error(err)
	}
}
