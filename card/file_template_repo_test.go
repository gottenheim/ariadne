package card_test

import (
	"testing"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/fs"
)

func TestFileTemplateRepository_GetTemplate(t *testing.T) {
	fakeFs, err := fs.NewFake([]fs.FakeEntry{
		fs.NewFakeEntry("/config/template", "source.cpp", `template source code file`),
		fs.NewFakeEntry("/config/template", "header.h", `template header file`),
	})

	if err != nil {
		t.Fatal(err)
	}

	repo := card.NewFileTemplateRepository(fakeFs, "/config/template")

	cardTemplate, err := repo.GetTemplate()

	if err != nil {
		t.Fatal(err)
	}

	sourceArtifact := cardTemplate.FindArtifactByName("source.cpp")

	if sourceArtifact == nil {
		t.Fatal("Source artifact is missing in card template")
	}

	if string(sourceArtifact.Content()) != "template source code file" {
		t.Fatal("Source artifact has unexpected content")
	}

	headerArtifact := cardTemplate.FindArtifactByName("header.h")

	if headerArtifact == nil {
		t.Fatal("Header artifact is missing in card template")
	}

	if string(headerArtifact.Content()) != "template header file" {
		t.Fatal("Header artifact has unexpected content")
	}
}
