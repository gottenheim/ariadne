package fs

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type FakeFileEntry struct {
	path     string
	name     string
	contents []byte
}

func NewFakeFileEntry(path string, name string, contents string) FakeFileEntry {
	return FakeFileEntry{
		path:     path,
		name:     name,
		contents: []byte(contents),
	}
}

func NewFakeFs(entries []FakeFileEntry) (afero.Fs, error) {
	fs := afero.NewMemMapFs()

	err := AddFakeFileEntries(fs, entries)

	if err != nil {
		return nil, err
	}

	return fs, nil
}

func AddFakeFileEntries(fs afero.Fs, entries []FakeFileEntry) error {
	for _, dirEntry := range entries {
		err := fs.MkdirAll(dirEntry.path, os.ModePerm)
		if err != nil {
			return errors.WithMessage(err,
				fmt.Sprintf("Failed to create directory %s in memory mapping filesystem",
					dirEntry.path))
		}

		filePath := path.Join(dirEntry.path, dirEntry.name)

		err = afero.WriteFile(fs, filePath, []byte(dirEntry.contents), os.ModePerm)

		if err != nil {
			return errors.WithMessage(err,
				fmt.Sprintf("Failed to create file %s in memory mapping filesystem",
					filePath))
		}
	}

	return nil
}
