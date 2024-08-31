package backend

import "github.com/bluemir/wikinote/internal/backend/files"

func (backend *Backend) FileSearch(query string) (*files.SearchResult, error) {
	//TODO query to pattern(regexp)
	return backend.files.Search(query)
}
