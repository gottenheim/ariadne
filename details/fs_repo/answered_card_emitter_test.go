package fs_repo_test

import (
	"testing"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/details/fs_repo"
	"github.com/gottenheim/ariadne/libraries/datetime"
	"github.com/gottenheim/ariadne/libraries/fs"
	"github.com/gottenheim/ariadne/libraries/pipeline"
)

func TestAnsweredCardEmitter_ShouldVisitAllCardSubdirectories(t *testing.T) {
	answeredCards := emitCards(t, "/home/user/books/cpp",
		fs.NewFakeFileEntry("/home/user/books/cpp/stroustrup/1", card.AnswerArtifactName, `fake answer`),
		fs.NewFakeFileEntry("/home/user/books/cpp/stroustrup/2", card.AnswerArtifactName, `fake answer`),
		fs.NewFakeFileEntry("/home/user/books/cpp/lippman/1", card.AnswerArtifactName, `fake answer`),
		fs.NewFakeFileEntry("/home/user/books/cpp/lippman/2", card.AnswerArtifactName, `fake answer`))

	assertCardFound(t, answeredCards, "/home/user/books/cpp/stroustrup", "1")
	assertCardFound(t, answeredCards, "/home/user/books/cpp/stroustrup", "2")
	assertCardFound(t, answeredCards, "/home/user/books/cpp/lippman", "1")
	assertCardFound(t, answeredCards, "/home/user/books/cpp/lippman", "2")
}

func TestAnsweredCardEmitter_ShouldNotVisitDirectoriesOutsideGivenDirectory(t *testing.T) {
	answeredCards := emitCards(t, "/home/user/books/cpp",
		fs.NewFakeFileEntry("/home/user/books/cpp/stroustrup/1", card.AnswerArtifactName, `fake answer`),
		fs.NewFakeFileEntry("/home/user/books/cpp/stroustrup/2", card.AnswerArtifactName, `fake answer`),
		fs.NewFakeFileEntry("/home/user/books/c#/in_the_nutshell/1", card.AnswerArtifactName, `fake answer`),
		fs.NewFakeFileEntry("/home/user/books/c#/in_the_nutshell/2", card.AnswerArtifactName, `fake answer`))

	if len(answeredCards) != 2 {
		t.Error("Two cards are expected to be found")
	}
	assertCardFound(t, answeredCards, "/home/user/books/cpp/stroustrup", "1")
	assertCardFound(t, answeredCards, "/home/user/books/cpp/stroustrup", "2")
}

func TestAnsweredCardEmitter_ShouldNotVisitCardsWithoutAnswers(t *testing.T) {
	answeredCards := emitCards(t, "/home/user/books/cpp",
		fs.NewFakeFileEntry("/home/user/books/cpp/stroustrup/1", card.AnswerArtifactName, `fake answer`),
		fs.NewFakeFileEntry("/home/user/books/cpp/stroustrup/2", "question.cpp", `source code file`))

	if len(answeredCards) != 1 {
		t.Error("One card is expected to be found")
	}
	assertCardFound(t, answeredCards, "/home/user/books/cpp/stroustrup", "1")
}

func TestAnsweredCardEmitter_ShouldCreateLearnCardActivity_IfNoActivitiesFileIsFound(t *testing.T) {
	answeredCards := emitCards(t, "/home/user/books/cpp",
		fs.NewFakeFileEntry("/home/user/books/cpp/stroustrup/1", card.AnswerArtifactName, `fake answer`))

	_, isLearnCard := answeredCards[0].Activities.(*card.LearnCardActivity)

	if !isLearnCard {
		t.Fatal("Learn card activity must be generated for directory without activities file")
	}
}

func TestAnsweredCardEmitter_ShouldReadExistingActivities_IfFileExists(t *testing.T) {
	fakeFs, err := fs.NewFakeFs([]fs.FakeFileEntry{})

	if err != nil {
		t.Fatal(err)
	}

	repo := fs_repo.NewFileCardRepository(fakeFs)

	c := card.CreateNew("/home/user/books/cpp/stroustrup",
		[]card.CardArtifact{
			card.NewCardArtifact("source.cpp", []byte("source code artifact")),
			card.NewCardArtifact("header.h", []byte("header artifact")),
		})

	c.StoreAnswer()

	c.SetActivities(card.GenerateActivityChain(card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToToday))

	err = repo.Save(c)

	if err != nil {
		t.Fatal(err)
	}

	p := pipeline.New()

	answeredCards := pipeline.NewItemCollector[card.BriefCard]()

	cardEmitter := pipeline.NewEmitter(p, fs_repo.NewAnsweredCardEmitter(fakeFs, repo, "/home/user/books/cpp"))

	pipeline.WithAcceptor[card.BriefCard](p, cardEmitter, answeredCards)

	err = p.SyncRun()

	if err != nil {
		t.Fatal(err)
	}

	answeredCard := answeredCards.Items[0]

	isScheduledToRemindToday, err := card.IsCardScheduledToRemindToday(datetime.NewFakeTimeSource(), answeredCard.Activities)

	if err != nil {
		t.Fatal(err)
	}

	if !isScheduledToRemindToday {
		t.Fatal("Card must be scheduled to remind today")
	}
}

func assertCardFound(t *testing.T, cards []card.BriefCard, section string, entry string) {
	for _, card := range cards {
		if card.Section == section && card.Entry == entry {
			return
		}
	}
	t.Fatalf("Card with section %s and entry %s was not found", section, entry)
}

func emitCards(t *testing.T, cardsDir string, files ...fs.FakeFileEntry) []card.BriefCard {
	fakeFs, err := fs.NewFakeFs(files)

	if err != nil {
		t.Fatal(err)
	}

	repo := fs_repo.NewFileCardRepository(fakeFs)

	p := pipeline.New()

	answeredCards := pipeline.NewItemCollector[card.BriefCard]()

	cardEmitter := pipeline.NewEmitter(p, fs_repo.NewAnsweredCardEmitter(fakeFs, repo, cardsDir))

	pipeline.WithAcceptor[card.BriefCard](p, cardEmitter, answeredCards)

	err = p.SyncRun()

	if err != nil {
		t.Fatal(err)
	}

	return answeredCards.Items
}
