package card_test

import (
	"testing"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/fs"
	"github.com/spf13/afero"
)

func TestFirstCard(t *testing.T) {
	fakeFs, err := fs.NewFake([]fs.FakeEntry{
		fs.NewFakeEntry("/config/template", "question.cpp", `template question`),
		fs.NewFakeEntry("/config/template", "header.h", `template question header`),
	})

	if err != nil {
		t.Fatal(err)
	}

	config := &card.Config{
		TemplateDir:    "/config/template",
		AnswerFileName: "answer.tgz",
	}

	card := card.New(fakeFs, config)
	_, err = card.CreateNew("/books/cpp")

	if err != nil {
		t.Fatal(err)
	}

	questionText, err := afero.ReadFile(fakeFs, "/books/cpp/1/question.cpp")

	if err != nil {
		t.Fatal(err)
	}

	if string(questionText) != "template question" {
		t.Error("Question texts don't match")
	}

	headerText, err := afero.ReadFile(fakeFs, "/books/cpp/1/header.h")

	if err != nil {
		t.Fatal(err)
	}

	if string(headerText) != "template question header" {
		t.Error("Header texts don't match")
	}
}

func TestSecondCard(t *testing.T) {
	fakeFs, err := fs.NewFake([]fs.FakeEntry{
		fs.NewFakeEntry("/config/template", "question.cpp", `template question`),
		fs.NewFakeEntry("/config/template", "header.h", `template question header`),
		fs.NewFakeEntry("/books/cpp/1", "answer.tgz", `answer number 1`),
		fs.NewFakeEntry("/books/cpp/1", "question.cpp", `question number 1`),
	})

	if err != nil {
		t.Fatal(err)
	}

	config := &card.Config{
		TemplateDir:    "/config/template",
		AnswerFileName: "answer.tgz",
	}

	card := card.New(fakeFs, config)
	_, err = card.CreateNew("/books/cpp")

	if err != nil {
		t.Fatal(err)
	}

	questionText, err := afero.ReadFile(fakeFs, "/books/cpp/2/question.cpp")

	if err != nil {
		t.Fatal(err)
	}

	if string(questionText) != "template question" {
		t.Error("Question texts don't match")
	}

	headerText, err := afero.ReadFile(fakeFs, "/books/cpp/2/header.h")

	if err != nil {
		t.Fatal(err)
	}

	if string(headerText) != "template question header" {
		t.Error("Header texts don't match")
	}
}
