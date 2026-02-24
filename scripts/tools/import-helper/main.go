package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	output := flag.String("o", "", "output file path (e.g. assets/src/js/index.js)")
	flag.Parse()

	if *output == "" || flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: import-helper -o <output-file> <source-dir>")
		os.Exit(1)
	}

	dir := flag.Arg(0)
	outputAbs, _ := filepath.Abs(*output)

	jsFiles := []string{}

	if err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".js") {
			return nil
		}

		abs, _ := filepath.Abs(path)
		if abs == outputAbs { // remove self reference
			return nil
		}

		path = strings.TrimPrefix(path, dir+"/")

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

	existing, err := os.ReadFile(*output)
	if err == nil && string(existing) == data {
		return // content unchanged, skip write to preserve timestamp
	}

	if err := os.WriteFile(*output, []byte(data), 0644); err != nil {
		panic(err)
	}
}
