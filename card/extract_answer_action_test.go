package card_test

import (
	"os"
	"testing"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/fs"
	"github.com/gottenheim/ariadne/test"
	"github.com/spf13/afero"
)

func TestExtractAnswerAction(t *testing.T) {
	fakeFs, err := fs.NewFake([]fs.FakeEntry{
		fs.NewFakeEntry("/home/user/books/cpp/1", "source.cpp", `old source code artifact`),
		fs.NewFakeEntry("/home/user/books/cpp/1", "header.h", `old header artifact`),
	})

	if err != nil {
		t.Fatal(err)
	}

	cardRepo := card.NewFileCardRepository(fakeFs, "/home/user")

	c, err := cardRepo.Get("books/cpp/1")

	if err != nil {
		t.Fatal(err)
	}

	c.CompressAnswer()

	err = cardRepo.Save(c)

	afero.WriteFile(fakeFs, "/home/user/books/cpp/1/source.cpp", []byte("new source code artifact"), os.ModePerm)
	afero.WriteFile(fakeFs, "/home/user/books/cpp/1/header.h", []byte("new header artifact"), os.ModePerm)

	extractAnswerAction := &card.ExtractCardAction{}
	err = extractAnswerAction.Run(fakeFs, "/home/user/", "/books/cpp/1")

	if err != nil {
		t.Fatal(err)
	}

	test.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/1/source.cpp", "old source code artifact")
	test.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/1/header.h", "old header artifact")
}
