package files

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

type SearchResultFile struct {
	FileName string
	Items    []SearchResultItem
}

type SearchResultItem struct {
	Line int    `json:"line"`
	Text string `json:"text"`
}
type SearchResult struct {
	Files []SearchResultFile
	Total int
}

func (fs *FileStore) Search(query string) (*SearchResult, error) {
	ret := SearchResult{}
	err := filepath.Walk(fs.wikipath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
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
		case ".md", ".txt", ".html":
			break
		default:
			return nil
		}

		items, err := fileSearch(path, query)
		if err != nil {
			return err
		}
		if len(items) != 0 {
			path, err := filepath.Rel(fs.wikipath, path)
			if err != nil {
				return err
			}
			ret.Files = append(ret.Files, SearchResultFile{
				FileName: "/" + path,
				Items:    items,
			})

			ret.Total += len(items)
		}
		return nil
	})
	return &ret, err
}
func fileSearch(path, query string) ([]SearchResultItem, error) {
	result := []SearchResultItem{}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	linenum := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		linenum++
		// TODO ignore case
		if strings.Contains(scanner.Text(), query) {
			result = append(result, SearchResultItem{
				Line: linenum,
				Text: scanner.Text(),
			})
		}
	}

	return result, scanner.Err()
}
