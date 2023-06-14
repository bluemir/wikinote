package assets

import "io/fs"

var (
	HTMLTemplates fs.FS
	Static        fs.FS
)

func InitFS(rootfs fs.FS) error {
	var err error
	Static, err = fs.Sub(rootfs, "build/static")
	if err != nil {
		return err
	}
	HTMLTemplates, err = fs.Sub(rootfs, "assets/html-templates")
	if err != nil {
		return err
	}

	return nil
}
