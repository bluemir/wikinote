package backend

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

type SearchResult struct {
	Line int
	Text string
}

func search(basepath string, query string) (interface{}, error) {

	result := map[string]([]SearchResult){}
	err := filepath.Walk(basepath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name()[0] == '.' {
			logrus.Debug(info.Name())
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		if info.Name()[0] == '.' {
			return nil
		}
		switch filepath.Ext(info.Name()) {
		case ".pdf", ".mp4", ".mp3":
			// if binary
			return nil
		case ".md":
			break
		default:
			return nil
		}

		res, err := fileSearch(path, query)
		if err != nil {
			return err
		}
		if len(res) != 0 {
			result[path[len(basepath):]] = res
		}
		return nil
	})
	return result, err
}
func fileSearch(path, query string) ([]SearchResult, error) {
	result := []SearchResult{}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	linenum := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		linenum++
		if strings.Contains(scanner.Text(), query) {
			result = append(result, SearchResult{
				Line: linenum,
				Text: scanner.Text(),
			})
		}
	}

	return result, scanner.Err()
}
