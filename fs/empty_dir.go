package fs

import (
	"github.com/spf13/afero"
	"path/filepath"
)

func EmptyDirectory(fs afero.Fs, dir string) error {
	d, err := fs.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = fs.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
