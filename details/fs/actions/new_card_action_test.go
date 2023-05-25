package card_test

import (
	"testing"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/fs"
	"github.com/gottenheim/ariadne/details/fs/card_repo"
)

func TestCreateNewCardInEmptyDirectory(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{
		fs.NewFakeFileEntry("/config/template", "source.cpp", `template source code file`),
		fs.NewFakeFileEntry("/config/template", "header.h", `template header file`),
	})

	if err != nil {
		t.Fatal(err)
	}

	templateRepo := card_repo.NewFileTemplateRepository(fakeFs, "/config/template")
	cardRepo := card_repo.NewFileCardRepository(fakeFs, "/home/user")

	newCardAction := &card.NewCardAction{}
	err = newCardAction.Run(templateRepo, cardRepo, []string{"books"})

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

	templateRepo := card_repo.NewFileTemplateRepository(fakeFs, "/config/template")
	cardRepo := card_repo.NewFileCardRepository(fakeFs, "/home/user")

	newCardAction := &card.NewCardAction{}
	err = newCardAction.Run(templateRepo, cardRepo, []string{"books"})

	if err != nil {
		t.Fatal(err)
	}

	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/2/source.cpp", "template source code file")
	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/2/header.h", "template header file")
}
