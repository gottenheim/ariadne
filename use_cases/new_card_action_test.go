package use_cases_test

import (
	"testing"

	"github.com/gottenheim/ariadne/infra/repo/fs_repo"
	"github.com/gottenheim/ariadne/libraries/fs"
	"github.com/gottenheim/ariadne/use_cases"
)

func TestCreateNewCardInEmptyDirectory(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{
		fs.NewFakeFileEntry("/config/template", "source.cpp", `template source code file`),
		fs.NewFakeFileEntry("/config/template", "header.h", `template header file`),
	})

	if err != nil {
		t.Fatal(err)
	}

	templateRepo := fs_repo.NewFileTemplateRepository(fakeFs, "/config/template")
	cardRepo := fs_repo.NewFileCardRepository(fakeFs)

	newCardAction := &use_cases.NewCard{}
	err = newCardAction.Run(templateRepo, cardRepo, "/home/user/books")

	if err != nil {
		t.Fatal(err)
	}

	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/1/source.cpp", "template source code file")
	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/1/header.h", "template header file")
}

func TestCreateNewCardInDirectoryWithCards(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{
		fs.NewFakeFileEntry("/config/template", "source.cpp", `template source code file`),
		fs.NewFakeFileEntry("/config/template", "header.h", `template header file`),
		fs.NewFakeFileEntry("/home/user/books/cpp/1", "source.cpp", `source code artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/1", "header.h", `header artifact`),
	})

	if err != nil {
		t.Fatal(err)
	}

	templateRepo := fs_repo.NewFileTemplateRepository(fakeFs, "/config/template")
	cardRepo := fs_repo.NewFileCardRepository(fakeFs)

	newCardAction := &use_cases.NewCard{}
	err = newCardAction.Run(templateRepo, cardRepo, "/home/user/books")

	if err != nil {
		t.Fatal(err)
	}

	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/2/source.cpp", "template source code file")
	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/2/header.h", "template header file")
}
