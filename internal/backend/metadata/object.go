package metadata

type ObjectStorageStore struct {
}

func (store *ObjectStorageStore) Take(path, key string) (string, error) {
	return "", ErrNotImplemented
}
func (store *ObjectStorageStore) Save(path, key, value string) error {
	return ErrNotImplemented
}
func (store *ObjectStorageStore) Delete(path, key string) error {
	return ErrNotImplemented
}
