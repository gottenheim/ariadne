package card_repo_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/config"
	"github.com/gottenheim/ariadne/details/fs"
	"github.com/gottenheim/ariadne/details/fs/card_repo"
)

func TestSavingFirstCardInEmptyRepository(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{})

	if err != nil {
		t.Fatal(err)
	}

	c := card.NewCard([]string{"books", "cpp"}, 0,
		[]card.CardArtifact{
			card.NewCardArtifact("source.cpp", []byte("source code artifact")),
			card.NewCardArtifact("header.h", []byte("header artifact")),
			card.NewCardArtifact("config.yml", []byte("config file artifact")),
		})

	repo := card_repo.NewFileCardRepository(fakeFs, "/home/user/")

	err = repo.Save(c)

	if err != nil {
		t.Fatal(err)
	}

	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/1/source.cpp", "source code artifact")
	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/1/header.h", "header artifact")
	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/1/config.yml", "config file artifact")
}

func TestSavingNewCardInExistingRepository(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{
		fs.NewFakeFileEntry("/home/user/books/cpp/1", "source.cpp", `1st source code artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/1", "header.h", `1st header artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/1", "config.yml", `1st config file`),
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

	repo := card_repo.NewFileCardRepository(fakeFs, "/home/user/")

	err = repo.Save(c)

	if err != nil {
		t.Fatal(err)
	}

	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/2/source.cpp", "2nd source code artifact")
	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/2/header.h", "2nd header artifact")
	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/2/config.yml", "2nd config file artifact")
}

func TestOverwritingCardInExistingRepository(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{
		fs.NewFakeFileEntry("/home/user/books/cpp/1", "source.cpp", `old source code artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/1", "header.h", `old header artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/1", "config.yml", `old config file`),
	})

	if err != nil {
		t.Fatal(err)
	}

	c := card.NewCard([]string{"books", "cpp"}, 1,
		[]card.CardArtifact{
			card.NewCardArtifact("source.cpp", []byte("new source code artifact")),
			card.NewCardArtifact("header.h", []byte("new header artifact")),
		})

	repo := card_repo.NewFileCardRepository(fakeFs, "/home/user")

	err = repo.Save(c)

	if err != nil {
		t.Fatal(err)
	}

	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/1/source.cpp", "new source code artifact")
	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/1/header.h", "new header artifact")
	fs.AssertFileDoesNotExists(t, fakeFs, "/home/user/books/cpp/1/config.yml")
}

func TestGetCardFromRepository(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{
		fs.NewFakeFileEntry("/home/user/books/cpp/2", "source.cpp", `source code artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/2", "header.h", `header artifact`),
	})

	if err != nil {
		t.Fatal(err)
	}

	repo := card_repo.NewFileCardRepository(fakeFs, "/home/user")

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
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{})

	if err != nil {
		t.Fatal(err)
	}

	repo := card_repo.NewFileCardRepository(fakeFs, "/home/user")

	c := card.NewCard([]string{"books", "cpp"}, 0,
		[]card.CardArtifact{
			card.NewCardArtifact("source.cpp", []byte("source code artifact")),
			card.NewCardArtifact("header.h", []byte("header artifact")),
		})

	err = repo.Save(c)

	if err != nil {
		t.Fatal(err)
	}

	progressFilePath := fmt.Sprintf("/home/user/books/cpp/1/%s", card_repo.ActivitiesFileName)

	fs.AssertFileDoesNotExists(t, fakeFs, progressFilePath)
	fs.AssertDirectoryFilesCount(t, fakeFs, filepath.Dir(progressFilePath), 2)
}

func TestSaveCreatesCardActivitiesFileIfStatusIsNotNew(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{})

	if err != nil {
		t.Fatal(err)
	}

	repo := card_repo.NewFileCardRepository(fakeFs, "/home/user")

	c := card.NewCard([]string{"books", "cpp"}, 0,
		[]card.CardArtifact{
			card.NewCardArtifact("source.cpp", []byte("source code artifact")),
			card.NewCardArtifact("header.h", []byte("header artifact")),
		})

	cardActivities := card.GenerateActivityChain(card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToYesterday|card.CardExecutedToday)

	c.SetActivities(cardActivities)

	err = repo.Save(c)

	if err != nil {
		t.Fatal(err)
	}

	activitiesBinary, err := card.SerializeCardActivityChain(cardActivities)
	if err != nil {
		t.Fatal(err)
	}

	activitiesFilePath := fmt.Sprintf("/home/user/books/cpp/1/%s", card_repo.ActivitiesFileName)

	fs.AssertFileExistsAndHasYamlContent(t, fakeFs, activitiesFilePath, string(activitiesBinary))
}

func TestReadCardActivities(t *testing.T) {
	initialActivities := card.GenerateActivityChain(card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToYesterday|card.CardExecutedToday)

	initialActivitiesBinary, err := card.SerializeCardActivityChain(initialActivities)
	if err != nil {
		t.Fatal(err)
	}

	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{
		fs.NewFakeFileEntry("/home/user/books/cpp/2", "source.cpp", `source code artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/2", "header.h", `header artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/2", card_repo.ActivitiesFileName, string(initialActivitiesBinary)),
	})

	if err != nil {
		t.Fatal(err)
	}

	repo := card_repo.NewFileCardRepository(fakeFs, "/home/user")

	c, err := repo.Get("/books/cpp/2")

	if err != nil {
		t.Fatal(err)
	}

	actualActivities := c.Activities()

	actualActivitiesBinary, err := card.SerializeCardActivityChain(actualActivities)
	if err != nil {
		t.Fatal(err)
	}

	config.AssertIdenticalYamlStrings(t, string(initialActivitiesBinary), string(actualActivitiesBinary))
}
