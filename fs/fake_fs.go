package fs

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type FakeEntry struct {
	path     string
	name     string
	contents []byte
}

func NewFakeEntry(path string, name string, contents string) FakeEntry {
	return FakeEntry{
		path:     path,
		name:     name,
		contents: []byte(contents),
	}
}

func NewFake(entries []FakeEntry) (afero.Fs, error) {
	fs := afero.NewMemMapFs()

	err := AddFakeEntries(fs, entries)

	if err != nil {
		return nil, err
	}

	return fs, nil
}

func AddFakeEntries(fs afero.Fs, entries []FakeEntry) error {
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
