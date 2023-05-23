package test

import (
	"testing"

	"github.com/spf13/afero"
)

func AssertFileExistsAndHasContent(t *testing.T, fs afero.Fs, path string, expectedContent string) {
	fileText, err := afero.ReadFile(fs, path)

	if err != nil {
		t.Errorf("File %s does not exist", path)
	}

	if string(fileText) != expectedContent {
		t.Errorf("File %s has unexpected content. Actual: %s, expected: %s", path, fileText, expectedContent)
	}
}

func AssertFileDoesNotExists(t *testing.T, fs afero.Fs, path string) {
	exists, err := afero.Exists(fs, path)

	if err != nil {
		t.Errorf("Unable to check if file %s exist or not", path)
	}

	if exists {
		t.Errorf("File %s exists, but it is not expected to be", path)
	}
}
