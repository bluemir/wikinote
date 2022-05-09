package backend

func (backend *Backend) FileRead(path string) ([]byte, error) {
	// TODO post hook?
	return backend.files.Read(path)
}

func (backend *Backend) FileWrite(path string, data []byte) error {
	data, err := backend.Plugin.TriggerFileWriteHook(path, data)
	if err != nil {
		return err
	}

	return backend.files.Write(path, data)
}

func (backend *Backend) FileDelete(path string) error {
	return backend.files.Delete(path)
}
