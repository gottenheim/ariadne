package card_test

import (
	"bytes"
	"testing"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/archive"
	"github.com/gottenheim/ariadne/details/fs"
	"github.com/gottenheim/ariadne/details/fs/card_repo"
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

	cardRepo := card_repo.NewFileCardRepository(fakeFs, "/home/user")

	comressAnswerAction := &card.CompressAnswerAction{}
	err = comressAnswerAction.Run(cardRepo, "/books/cpp/1")

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
