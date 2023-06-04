package fs_repo_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/details/fs_repo"
	"github.com/gottenheim/ariadne/libraries/config"
	"github.com/gottenheim/ariadne/libraries/fs"
)

func TestSavingFirstCardInEmptyRepository(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{})

	if err != nil {
		t.Fatal(err)
	}

	c := card.CreateNew("/home/user/books/cpp",
		[]card.CardArtifact{
			card.NewCardArtifact("source.cpp", []byte("source code artifact")),
			card.NewCardArtifact("header.h", []byte("header artifact")),
			card.NewCardArtifact("config.yml", []byte("config file artifact")),
		})

	repo := fs_repo.NewFileCardRepository(fakeFs)

	err = repo.Save(c)

	if err != nil {
		t.Fatal(err)
	}

	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/01/source.cpp", "source code artifact")
	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/01/header.h", "header artifact")
	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/01/config.yml", "config file artifact")
}

func TestSavingNewCardInExistingRepository(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{
		fs.NewFakeFileEntry("/home/user/books/cpp/01", "source.cpp", `1st source code artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/01", "header.h", `1st header artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/01", "config.yml", `1st config file`),
	})

	if err != nil {
		t.Fatal(err)
	}

	c := card.CreateNew("/home/user/books/cpp",
		[]card.CardArtifact{
			card.NewCardArtifact("source.cpp", []byte("2nd source code artifact")),
			card.NewCardArtifact("header.h", []byte("2nd header artifact")),
			card.NewCardArtifact("config.yml", []byte("2nd config file artifact")),
		})

	repo := fs_repo.NewFileCardRepository(fakeFs)

	err = repo.Save(c)

	if err != nil {
		t.Fatal(err)
	}

	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/02/source.cpp", "2nd source code artifact")
	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/02/header.h", "2nd header artifact")
	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/02/config.yml", "2nd config file artifact")
}

func TestOverwritingCardInExistingRepository(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{
		fs.NewFakeFileEntry("/home/user/books/cpp/01", "source.cpp", `old source code artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/01", "header.h", `old header artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/01", "config.yml", `old config file`),
	})

	if err != nil {
		t.Fatal(err)
	}

	c := card.FromExisting("/home/user/books/cpp", "01",
		[]card.CardArtifact{
			card.NewCardArtifact("source.cpp", []byte("new source code artifact")),
			card.NewCardArtifact("header.h", []byte("new header artifact")),
		}, card.CreateLearnCardActivity())

	repo := fs_repo.NewFileCardRepository(fakeFs)

	err = repo.Save(c)

	if err != nil {
		t.Fatal(err)
	}

	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/01/source.cpp", "new source code artifact")
	fs.AssertFileExistsAndHasContent(t, fakeFs, "/home/user/books/cpp/01/header.h", "new header artifact")
	fs.AssertFileDoesNotExists(t, fakeFs, "/home/user/books/cpp/01/config.yml")
}

func TestGetCardFromRepository(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{
		fs.NewFakeFileEntry("/home/user/books/cpp/02", "source.cpp", `source code artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/02", "header.h", `header artifact`),
	})

	if err != nil {
		t.Fatal(err)
	}

	repo := fs_repo.NewFileCardRepository(fakeFs)

	c, err := repo.Get("/home/user/books/cpp", "02")

	if err != nil {
		t.Fatal(err)
	}

	if c.Section() != "/home/user/books/cpp" {
		t.Error("Loaded card has unexpected section")
	}

	if c.Entry() != "02" {
		t.Error("Loaded card has unexpected entry")
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

func TestSkippingActivitiesFileDuringArtifactReading(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{
		fs.NewFakeFileEntry("/home/user/books/cpp/02", "source.cpp", `source code artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/02", "header.h", `header artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/02", ".activities", `ActivityType: learn`),
	})

	if err != nil {
		t.Fatal(err)
	}

	repo := fs_repo.NewFileCardRepository(fakeFs)

	c, err := repo.Get("/home/user/books/cpp", "02")

	if err != nil {
		t.Fatal(err)
	}

	activitiesFile := c.FindArtifactByName(".activities")
	if activitiesFile != nil {
		t.Fatal("Activities file must not be added to artifact collection")
	}
}

func TestSkipNewCardProgressDuringSaving(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{})

	if err != nil {
		t.Fatal(err)
	}

	repo := fs_repo.NewFileCardRepository(fakeFs)

	c := card.CreateNew("/home/user/books/cpp",
		[]card.CardArtifact{
			card.NewCardArtifact("source.cpp", []byte("source code artifact")),
			card.NewCardArtifact("header.h", []byte("header artifact")),
		})

	err = repo.Save(c)

	if err != nil {
		t.Fatal(err)
	}

	progressFilePath := fmt.Sprintf("/home/user/books/cpp/01/%s", fs_repo.ActivitiesFileName)

	fs.AssertFileDoesNotExists(t, fakeFs, progressFilePath)
	fs.AssertDirectoryFilesCount(t, fakeFs, filepath.Dir(progressFilePath), 2)
}

func TestSaveCreatesCardActivitiesFileIfStatusIsNotNew(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{})

	if err != nil {
		t.Fatal(err)
	}

	repo := fs_repo.NewFileCardRepository(fakeFs)

	c := card.CreateNew("/home/user/books/cpp",
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

	activitiesFilePath := fmt.Sprintf("/home/user/books/cpp/01/%s", fs_repo.ActivitiesFileName)

	fs.AssertFileExistsAndHasYamlContent(t, fakeFs, activitiesFilePath, string(activitiesBinary))
}

func TestReadCardActivities(t *testing.T) {
	initialActivities := card.GenerateActivityChain(card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToYesterday|card.CardExecutedToday)

	initialActivitiesBinary, err := card.SerializeCardActivityChain(initialActivities)
	if err != nil {
		t.Fatal(err)
	}

	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{
		fs.NewFakeFileEntry("/home/user/books/cpp/02", "source.cpp", `source code artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/02", "header.h", `header artifact`),
		fs.NewFakeFileEntry("/home/user/books/cpp/02", fs_repo.ActivitiesFileName, string(initialActivitiesBinary)),
	})

	if err != nil {
		t.Fatal(err)
	}

	repo := fs_repo.NewFileCardRepository(fakeFs)

	c, err := repo.Get("/home/user/books/cpp", "02")

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
