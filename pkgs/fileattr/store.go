package fileattr

import (
	"strings"

	"github.com/jinzhu/gorm"
)

// store.Where(fileAttr.Options{NameSpace: "plugin.bluemir.me"}).SortBy(VALUE, DESC).Find()
// store.Where(fileAttr.Options{Path:"front-page.md"}).SortBy(KEY, DESC).Find()
type store struct {
	db *gorm.DB
}
type Options struct {
	Path      string
	Namespace string
	Key       string
}
type FileAttrEntity struct {
	gorm.Model
	FileAttr
}
type storeClause struct {
	*store
	options        *Options
	orderType      OrderType
	orderDirection OrderDirection
	limit          int
}

func NewStore(db *gorm.DB) (Store, error) {
	if err := db.AutoMigrate(&FileAttrEntity{}).Error; err != nil {
		return nil, err
	}
	return &store{db}, nil
}

func (s *store) Where(opt *Options) WhereClause {
	return (&storeClause{store: s}).Where(opt)
}
func (s *store) SortBy(t OrderType, d OrderDirection) SortByClause {
	return (&storeClause{store: s}).SortBy(t, d)
}
func (s *store) Limit(limit int) LimitClause {
	return (&storeClause{store: s}).Limit(limit)
}
func (s *store) Find() ([]FileAttr, error) {
	return (&storeClause{store: s}).Find()
}
func (s *store) Path(path string) PathClause {
	return &pathClause{store: s, path: path}
}

func (s *storeClause) Where(opt *Options) WhereClause {
	s.options = opt
	return s
}
func (s *storeClause) SortBy(t OrderType, d OrderDirection) SortByClause {
	s.orderType = t
	s.orderDirection = d
	return s
}
func (s *storeClause) Limit(limit int) LimitClause {
	s.limit = limit
	return s
}
func (s *storeClause) Find() ([]FileAttr, error) {
	attrs := []FileAttrEntity{}

	err := s.db.Where(&FileAttrEntity{
		FileAttr: FileAttr{
			Path:      s.options.Path,
			Namespace: s.options.Namespace,
			Key:       s.options.Key,
		},
	}).Find(&attrs).Error
	if err != nil {
		return nil, err
	}

	result := []FileAttr{}
	for _, attr := range attrs {
		result = append(result, attr.FileAttr)
	}

	return result, nil
}
func (s *storeClause) Get() (*FileAttr, error) {
	attr := &FileAttrEntity{}

	err := s.db.Where(&FileAttrEntity{
		FileAttr: FileAttr{
			Path:      s.options.Path,
			Namespace: s.options.Namespace,
			Key:       s.options.Key,
		},
	}).First(attr).Error
	if err != nil {
		return nil, err
	}

	return &attr.FileAttr, nil
}

type pathClause struct {
	store *store
	path  string
}

func (s *pathClause) Get(namespaceKey string) (string, error) {
	namespace, k := attrKeySplit(namespaceKey)

	attr, err := (&storeClause{store: s.store}).Where(&Options{
		Path:      s.path,
		Namespace: namespace,
		Key:       k,
	}).Get()
	if err != nil {
		return "", err
	}
	return attr.Value, nil
}
func (s *pathClause) Set(namespaceKey, value string) error {
	namespace, k := attrKeySplit(namespaceKey)

	attr := &FileAttrEntity{FileAttr: FileAttr{Namespace: namespace, Key: k, Value: value, Path: s.path}}
	return s.store.db.Where(&FileAttrEntity{
		FileAttr: FileAttr{
			Path:      s.path,
			Namespace: namespace,
			Key:       k,
		},
	}).Assign(attr).FirstOrCreate(attr).Error
}
func (s *pathClause) All(namespace string) (map[string]string, error) {
	attrs, err := (&storeClause{store: s.store}).Where(&Options{
		Path:      s.path,
		Namespace: namespace,
	}).Find()
	if err != nil {
		return nil, err
	}

	kv := map[string]string{}
	for _, attr := range attrs {
		kv[attr.Namespace+"/"+attr.Key] = attr.Value
	}
	return kv, nil
}
func attrKeySplit(key string) (string, string) {
	arr := strings.SplitN(key, "/", 2)
	return arr[0], arr[1]
}
