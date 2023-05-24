package card_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/fs"
	"github.com/gottenheim/ariadne/test"
)

func TestSavingFirstCardInEmptyRepository(t *testing.T) {
	fakeFs, err := fs.NewFake([]fs.FakeEntry{})

	if err != nil {
		t.Fatal(err)
	}

	c := card.NewCard([]string{"books", "cpp"}, 0,
		[]card.CardArtifact{
			card.NewCardArtifact("source.cpp", []byte("source code artifact")),
			card.NewCardArtifact("header.h", []byte("header artifact")),
			card.NewCardArtifact("config.yml", []byte("config file artifact")),
		})

	repo := card.NewFileCardRepository(fakeFs, "/home/user/")

	err = repo.Save(c)

	if err != nil {
		t.Fatal(err)
	}

	test.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/1/source.cpp", "source code artifact")
	test.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/1/header.h", "header artifact")
	test.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/1/config.yml", "config file artifact")
}

func TestSavingNewCardInExistingRepository(t *testing.T) {
	fakeFs, err := fs.NewFake([]fs.FakeEntry{
		fs.NewFakeEntry("/home/user/books/cpp/1", "source.cpp", `1st source code artifact`),
		fs.NewFakeEntry("/home/user/books/cpp/1", "header.h", `1st header artifact`),
		fs.NewFakeEntry("/home/user/books/cpp/1", "config.yml", `1st config file`),
	})

	if err != nil {
		t.Fatal(err)
	}

	c := card.NewCard([]string{"books", "cpp"}, 0,
		[]card.CardArtifact{
			card.NewCardArtifact("source.cpp", []byte("2nd source code artifact")),
			card.NewCardArtifact("header.h", []byte("2nd header artifact")),
			card.NewCardArtifact("config.yml", []byte("2nd config file artifact")),
		})

	repo := card.NewFileCardRepository(fakeFs, "/home/user/")

	err = repo.Save(c)

	if err != nil {
		t.Fatal(err)
	}

	test.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/2/source.cpp", "2nd source code artifact")
	test.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/2/header.h", "2nd header artifact")
	test.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/2/config.yml", "2nd config file artifact")
}

func TestOverwritingCardInExistingRepository(t *testing.T) {
	fakeFs, err := fs.NewFake([]fs.FakeEntry{
		fs.NewFakeEntry("/home/user/books/cpp/1", "source.cpp", `old source code artifact`),
		fs.NewFakeEntry("/home/user/books/cpp/1", "header.h", `old header artifact`),
		fs.NewFakeEntry("/home/user/books/cpp/1", "config.yml", `old config file`),
	})

	if err != nil {
		t.Fatal(err)
	}

	c := card.NewCard([]string{"books", "cpp"}, 1,
		[]card.CardArtifact{
			card.NewCardArtifact("source.cpp", []byte("new source code artifact")),
			card.NewCardArtifact("header.h", []byte("new header artifact")),
		})

	repo := card.NewFileCardRepository(fakeFs, "/home/user")

	err = repo.Save(c)

	if err != nil {
		t.Fatal(err)
	}

	test.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/1/source.cpp", "new source code artifact")
	test.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/1/header.h", "new header artifact")
	test.AssertFileDoesNotExists(t, fakeFs, "/home/user/books/cpp/1/config.yml")
}

func TestGetCardFromRepository(t *testing.T) {
	fakeFs, err := fs.NewFake([]fs.FakeEntry{
		fs.NewFakeEntry("/home/user/books/cpp/2", "source.cpp", `source code artifact`),
		fs.NewFakeEntry("/home/user/books/cpp/2", "header.h", `header artifact`),
	})

	if err != nil {
		t.Fatal(err)
	}

	repo := card.NewFileCardRepository(fakeFs, "/home/user")

	c, err := repo.Get("/books/cpp/2")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(c.Sections(), []string{"books", "cpp"}) {
		t.Error("Loaded card has unexpected sections")
	}

	if c.OrderNumber() != 2 {
		t.Error("Loaded card has unexpected order number")
	}

	sourceCode := c.FindArtifactByName("source.cpp")
	if sourceCode == nil {
		t.Fatal("Loaded card doesn't have source code artifact")
	}

	if string(sourceCode.Content()) != "source code artifact" {
		t.Error("Loaded card has unexpected source code artifact content")
	}

	headerFile := c.FindArtifactByName("header.h")
	if headerFile == nil {
		t.Fatal("Loaded card doesn't have header artifact")
	}

	if string(headerFile.Content()) != "header artifact" {
		t.Error("Loaded card has unexpected header artifact content")
	}
}

func TestSkipNewCardProgressDuringSaving(t *testing.T) {
	fakeFs, err := fs.NewFake([]fs.FakeEntry{})

	if err != nil {
		t.Fatal(err)
	}

	repo := card.NewFileCardRepository(fakeFs, "/home/user")

	c := card.NewCard([]string{"books", "cpp"}, 0,
		[]card.CardArtifact{
			card.NewCardArtifact("source.cpp", []byte("source code artifact")),
			card.NewCardArtifact("header.h", []byte("header artifact")),
		})

	err = repo.Save(c)

	if err != nil {
		t.Fatal(err)
	}

	progressFilePath := fmt.Sprintf("/home/user/books/cpp/1/%s", card.ProgressFileName)

	test.AssertFileDoesNotExists(t, fakeFs, progressFilePath)
	test.AssertDirectoryFilesCount(t, fakeFs, filepath.Dir(progressFilePath), 2)
}

func TestSaveCreatesCardProgressFileIfStatusIsNotNew(t *testing.T) {
	fakeFs, err := fs.NewFake([]fs.FakeEntry{})

	if err != nil {
		t.Fatal(err)
	}

	repo := card.NewFileCardRepository(fakeFs, "/home/user")

	c := card.NewCard([]string{"books", "cpp"}, 0,
		[]card.CardArtifact{
			card.NewCardArtifact("source.cpp", []byte("source code artifact")),
			card.NewCardArtifact("header.h", []byte("header artifact")),
		})

	c.SetProgress(card.ScheduledCard())

	err = repo.Save(c)

	if err != nil {
		t.Fatal(err)
	}

	progressFilePath := fmt.Sprintf("/home/user/books/cpp/1/%s", card.ProgressFileName)

	test.AssertFileExistsAndHasContent(t, fakeFs, progressFilePath, `Status: Scheduled`)
}

func TestReadCardProgress(t *testing.T) {
	fakeFs, err := fs.NewFake([]fs.FakeEntry{
		fs.NewFakeEntry("/home/user/books/cpp/2", "source.cpp", `source code artifact`),
		fs.NewFakeEntry("/home/user/books/cpp/2", "header.h", `header artifact`),
		fs.NewFakeEntry("/home/user/books/cpp/2", card.ProgressFileName, `Status: Scheduled`),
	})

	if err != nil {
		t.Fatal(err)
	}

	repo := card.NewFileCardRepository(fakeFs, "/home/user")

	c, err := repo.Get("/books/cpp/2")

	if err != nil {
		t.Fatal(err)
	}

	if !c.Progress().IsScheduled() {
		t.Fatal("Card progress status should be Scheduled")
	}
}
