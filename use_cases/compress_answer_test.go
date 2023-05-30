package use_cases_test

import (
	"bytes"
	"testing"

	"github.com/gottenheim/ariadne/infra/repo/fs_repo"
	"github.com/gottenheim/ariadne/libraries/archive"
	"github.com/gottenheim/ariadne/libraries/fs"
	"github.com/gottenheim/ariadne/use_cases"
	"github.com/spf13/afero"
)

func TestCompressAnswerAction(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{
		fs.NewFakeFileEntry("/home/user/books/cpp/1", "source.cpp", `source code artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/1", "header.h", `header artifact`),
	})

	if err != nil {
		t.Fatal(err)
	}

	cardRepo := fs_repo.NewFileCardRepository(fakeFs, "/home/user/books/cpp")

	comressAnswerAction := &use_cases.CompressAnswerAction{}
	err = comressAnswerAction.Run(cardRepo, 1)

	if err != nil {
		t.Fatal(err)
	}

	fileText, err := afero.ReadFile(fakeFs, "/home/user/books/cpp/1/answer.tgz")

	if err != nil {
		t.Fatal("Answer file was not generated")
	}

	files, err := archive.GetFiles(bytes.NewReader(fileText))

	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 2 {
		t.Error("Archive is expected to contain two files")
	}

	if string(files["source.cpp"]) != "source code artifact" {
		t.Error("Archive contains corrupted source file")
	}

	if string(files["header.h"]) != "header artifact" {
		t.Error("Archive contains corrupted source file")
	}
}
