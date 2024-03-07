//go:build !embed

package main

import (
	"os"

	"github.com/bluemir/wikinote/internal/assets"
)

func init() {
	// default, when no embed.
	if err := assets.InitFS(os.DirFS("./")); err != nil {
		panic(err)
	}
}
