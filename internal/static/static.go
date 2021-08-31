package static

import (
	rice "github.com/GeertJohan/go.rice"
)

var (
	Files         = rice.MustFindBox("../../build/static")
	HTMLTemplates = rice.MustFindBox("../../web/html-templates")
)
