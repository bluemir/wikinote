package backend

import (
	"github.com/bluemir/wikinote/internal/backend/attr"
)

func (backend *Backend) AttrFind(attr *attr.Attribute) ([]attr.Attribute, error) {
	return backend.attr.Find(attr)
}
func (backend *Backend) AttrSave(attr *attr.Attribute) error {
	return backend.attr.Save(attr)
}

func (backend *Backend) AttrDelete(attr *attr.Attribute) error {
	return backend.attr.Delete(attr)
}
