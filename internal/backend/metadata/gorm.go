package metadata

type GormStore struct {
}

func (store *GormStore) Take(path, key string) (string, error) {
	return "", ErrNotImplemented
}
func (store *GormStore) Save(path, key, value string) error {
	return ErrNotImplemented
}
func (store *GormStore) Delete(path, key string) error {
	return ErrNotImplemented
}
