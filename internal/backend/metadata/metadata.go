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
	Gorm *GormStoreConfig
}

func New(conf *Config) (Store, error) {
	//validation
	if count(
		conf.File != nil,
		conf.ObjectStorage != nil,
		conf.Gorm != nil,
	) > 1 {
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
		if err := conf.Gorm.DB.AutoMigrate(&GormEntry{}); err != nil {
			return nil, err
		}
		return &GormStore{conf.Gorm.DB}, nil
	default:
		return nil, errors.New("there is no option")
	}
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
