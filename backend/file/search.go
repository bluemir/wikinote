package file

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

func (m *manager) Search(query string) (interface{}, error) {
	//TODO query to pattern(regexp)
	result := map[string]([]SearchResult){}
	err := filepath.Walk(m.basepath, func(path string, info os.FileInfo, err error) error {
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

		res, err := fileSearch(path, query)
		if err != nil {
			return err
		}
		if len(res) != 0 {
			result[path[len(m.basepath):]] = res
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
