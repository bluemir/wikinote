package metadata

import (
	"errors"
)

type Store interface {
	Take(path, key string) (string, error)
	Save(path, key, value string) error
	Delete(path, key string) error
}

type Config struct {
	File          *FileStoreConfig
	ObjectStorage *struct {
		Prefix string
	}
	Gorm *struct {
		DBPath string
	}
}

func New(conf *Config) (Store, error) {
	//validation
	switch count(
		conf.File != nil,
		conf.ObjectStorage != nil,
		conf.Gorm != nil,
	) {
	case 0:
		return nil, errors.New("there is no option")
	case 1:
		break
	default:
		return nil, errors.New("multiple option exist")
	}

	switch {
	case conf.File != nil:
		if conf.File.Path == "" {
			return nil, errors.New("metadata config invaild. file.path is empty")
		}
		return &FileStore{conf.File}, nil
	case conf.ObjectStorage != nil:
		return &ObjectStorageStore{}, nil
	case conf.Gorm != nil:
		return &GormStore{}, nil
	}

	return nil, nil
}

func count(bs ...bool) int {
	i := 0
	for _, b := range bs {
		if b {
			i++
		}
	}
	return i
}
