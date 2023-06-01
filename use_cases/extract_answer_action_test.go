package use_cases_test

import (
	"os"
	"testing"

	"github.com/gottenheim/ariadne/infra/fs/fs_repo"
	"github.com/gottenheim/ariadne/libraries/fs"
	"github.com/gottenheim/ariadne/use_cases"
	"github.com/spf13/afero"
)

func TestExtractAnswerAction(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{
		fs.NewFakeFileEntry("/home/user/books/cpp/1", "source.cpp", `old source code artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/1", "header.h", `old header artifact`),
	})

	if err != nil {
		t.Fatal(err)
	}

	cardRepo := fs_repo.NewFileCardRepository(fakeFs)

	c, err := cardRepo.Get("/home/user/books/cpp", "1")

	if err != nil {
		t.Fatal(err)
	}

	c.StoreAnswer()

	err = cardRepo.Save(c)

	afero.WriteFile(fakeFs, "/home/user/books/cpp/1/source.cpp", []byte("new source code artifact"), os.ModePerm)
	afero.WriteFile(fakeFs, "/home/user/books/cpp/1/header.h", []byte("new header artifact"), os.ModePerm)

	extractAnswerAction := &use_cases.ExtractCard{}
	err = extractAnswerAction.Run(cardRepo, "/home/user/books/cpp", "1")

	if err != nil {
		t.Fatal(err)
	}

	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/1/source.cpp", "old source code artifact")
	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/1/header.h", "old header artifact")
}
