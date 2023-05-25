package fs

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gottenheim/ariadne/details/config"
	"github.com/spf13/afero"
)

func AssertFileExistsAndHasContent(t *testing.T, fs afero.Fs, path string, expectedContent string) {
	fileText, err := afero.ReadFile(fs, path)

	if err != nil {
		t.Errorf("File %s does not exist", path)
	}

	if strings.TrimSpace(string(fileText)) != expectedContent {
		t.Errorf("File %s has unexpected content. Actual: %s, expected: %s", path, fileText, expectedContent)
	}
}

func AssertFileExistsAndHasYamlContent(t *testing.T, fs afero.Fs, path string, expectedContent string) {
	fileText, err := afero.ReadFile(fs, path)

	if err != nil {
		t.Errorf("File %s does not exist", path)
	}

	config.AssertIdenticalYamlStrings(t, string(fileText), expectedContent)
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

func AssertDirectoryFilesCount(t *testing.T, fs afero.Fs, path string, expectedCount int) {
	filesCount := 0
	err := afero.Walk(fs, path, func(filePath string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			filesCount++
		}
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if filesCount != expectedCount {
		t.Error(fmt.Sprintf("Directory %s is expected to have %d files, but have %d", path, expectedCount, filesCount))
	}
}
