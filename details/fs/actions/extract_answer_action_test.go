package card_test

import (
	"os"
	"testing"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/fs"
	"github.com/gottenheim/ariadne/details/fs/card_repo"
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

	cardRepo := card_repo.NewFileCardRepository(fakeFs, "/home/user/books/cpp")

	c, err := cardRepo.Get(1)

	if err != nil {
		t.Fatal(err)
	}

	c.CompressAnswer()

	err = cardRepo.Save(c)

	afero.WriteFile(fakeFs, "/home/user/books/cpp/1/source.cpp", []byte("new source code artifact"), os.ModePerm)
	afero.WriteFile(fakeFs, "/home/user/books/cpp/1/header.h", []byte("new header artifact"), os.ModePerm)

	extractAnswerAction := &card.ExtractCardAction{}
	err = extractAnswerAction.Run(cardRepo, 1)

	if err != nil {
		t.Fatal(err)
	}

	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/1/source.cpp", "old source code artifact")
	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/1/header.h", "old header artifact")
}
