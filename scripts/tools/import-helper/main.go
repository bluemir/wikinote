package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	jsFiles := []string{}

	if err := filepath.Walk("assets/js", func(path string, info fs.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".js") {
			return nil
		}

		path = strings.TrimPrefix(path, "assets/js/")

		if path == "index.js" { // remove self reference
			return nil
		}

		jsFiles = append(jsFiles, path)

		return nil
	}); err != nil {
		panic(err)
	}

	lines := []string{}

	for _, file := range jsFiles {
		lines = append(lines, fmt.Sprintf(`import "./%s";`, file))
	}

	data := strings.Join(lines, "\n")

	if err := os.WriteFile("assets/js/index.js", []byte(data), 0644); err != nil {
		panic(err)
	}
}
