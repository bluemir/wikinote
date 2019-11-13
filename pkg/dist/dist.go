package dist

import (
	rice "github.com/GeertJohan/go.rice"
)

var (
	Files     = rice.MustFindBox("../../build/dist")
	Templates = rice.MustFindBox("../../app/html")
)
