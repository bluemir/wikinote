package fileattr

import (
	"strings"

	"github.com/pkg/errors"
)

// front-page.md|plugin.bluemir.me/last-changed
// |plugin.bluemir.me/
// /front-page.md|/
// /aa/info/test|some.thing.com/info/page/limit-less
//

func parseQuery(query string) (*Options, error) {
	arr := strings.SplitN(query, "|", 2)
	if len(arr) != 2 {
		return nil, errors.Errorf("query parse failed, '|' not found")
	}

	path := arr[0]

	arr2 := strings.SplitN(arr[1], "/", 2)
	if len(arr2) != 2 {
		return nil, errors.Errorf("query parse failed, '/' not found")
	}

	namespace, key := arr2[0], arr2[1]

	return &Options{
		Path:      path,
		Namespace: namespace,
		Key:       key,
	}, nil
}
