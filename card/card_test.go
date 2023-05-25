package card_test

import (
	"bytes"
	"testing"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/archive"
)

func TestCompressAnswerArtifacts(t *testing.T) {
	c := card.NewCard([]string{}, 0, []card.CardArtifact{
		card.NewCardArtifact("source.cpp", []byte("source file contents")),
		card.NewCardArtifact("header.h", []byte("header file contents")),
	})

	err := c.CompressAnswer()

	if err != nil {
		t.Fatal("Failed to compress card artifacts")
	}

	artifact := c.FindArtifactByName(card.AnswerArtifactName)

	if artifact == nil {
		t.Fatal("Artifacts weren't compressed")
	}

	files, err := archive.GetFiles(bytes.NewReader(artifact.Content()))

	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 2 {
		t.Error("Archive is expected to contain two files")
	}

	if string(files["source.cpp"]) != "source file contents" {
		t.Error("Archive contains corrupted source file")
	}

	if string(files["header.h"]) != "header file contents" {
		t.Error("Archive contains corrupted source file")
	}
}

func TestExtractAnswerArtifacts(t *testing.T) {
	card1 := card.NewCard([]string{}, 0, []card.CardArtifact{
		card.NewCardArtifact("source.cpp", []byte("old source file contents")),
		card.NewCardArtifact("header.h", []byte("old header file contents")),
	})

	err := card1.CompressAnswer()

	if err != nil {
		t.Fatal("Failed to compress card artifacts")
	}

	card2 := card.NewCard([]string{}, 0, []card.CardArtifact{
		card.NewCardArtifact(card.AnswerArtifactName, card1.FindAnswerArtifact().Content()),
		card.NewCardArtifact("source.cpp", []byte("new source file contents")),
		card.NewCardArtifact("header.h", []byte("new header file contents")),
	})

	card2.ExtractAnswer()

	if string(card2.FindArtifactByName("source.cpp").Content()) != "old source file contents" {
		t.Error("Source code artifact has unexpected contents")
	}

	if string(card2.FindArtifactByName("header.h").Content()) != "old header file contents" {
		t.Error("Source header artifact has unexpected contents")
	}
}
