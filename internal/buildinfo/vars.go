package buildinfo

import "github.com/bluemir/wikinote/internal/static"

var (
	Version   string
	AppName   string
	BuildTime string = static.Files.MustString(".time")
)
