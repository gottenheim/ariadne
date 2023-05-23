package card_test

import (
	"testing"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/fs"
	"github.com/gottenheim/ariadne/test"
)

func TestCreateNewCardInEmptyDirectory(t *testing.T) {
	fakeFs, err := fs.NewFake([]fs.FakeEntry{
		fs.NewFakeEntry("/config/template", "source.cpp", `template source code file`),
		fs.NewFakeEntry("/config/template", "header.h", `template header file`),
	})

	if err != nil {
		t.Fatal(err)
	}

	newCardAction := &card.NewCardAction{}
	err = newCardAction.Run(fakeFs, "/home/user/", "/books", "/config/template")

	if err != nil {
		t.Fatal(err)
	}

	test.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/1/source.cpp", "template source code file")
	test.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/1/header.h", "template header file")
}

func TestCreateNewCardInDirectoryWithCards(t *testing.T) {
	fakeFs, err := fs.NewFake([]fs.FakeEntry{
		fs.NewFakeEntry("/config/template", "source.cpp", `template source code file`),
		fs.NewFakeEntry("/config/template", "header.h", `template header file`),
		fs.NewFakeEntry("/home/user/books/cpp/1", "source.cpp", `source code artifact`),
		fs.NewFakeEntry("/home/user/books/cpp/1", "header.h", `header artifact`),
	})

	if err != nil {
		t.Fatal(err)
	}

	newCardAction := &card.NewCardAction{}
	err = newCardAction.Run(fakeFs, "/home/user/", "/books", "/config/template")

	if err != nil {
		t.Fatal(err)
	}

	test.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/2/source.cpp", "template source code file")
	test.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/2/header.h", "template header file")
}
