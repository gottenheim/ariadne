package card_test

import (
	"bytes"
	"testing"

	"github.com/gottenheim/ariadne/archive"
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
	_, err = card.CreateCard("/books/cpp")

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
	_, err = card.CreateCard("/books/cpp")

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

func TestPackAnswer(t *testing.T) {
	fakeFs, err := fs.NewFake([]fs.FakeEntry{
		fs.NewFakeEntry("/books/cpp/1", "question.cpp", `question source file`),
		fs.NewFakeEntry("/books/cpp/1", "question.h", `question header file`),
	})

	if err != nil {
		t.Fatal(err)
	}

	config := &card.Config{
		AnswerFileName: "answer.tgz",
	}

	card := card.New(fakeFs, config)
	err = card.PackAnswer("/books/cpp/1")

	if err != nil {
		t.Fatal(err)
	}

	answerFileContents, err := afero.ReadFile(fakeFs, "/books/cpp/1/answer.tgz")

	if err != nil {
		t.Fatal(err)
	}

	files, err := archive.GetFiles(bytes.NewReader(answerFileContents))

	if err != nil {
		t.Fatal(err)
	}

	if files["question.cpp"] != "question source file" {
		t.Error("question.cpp was not compressed into answer file")
	}

	if files["question.h"] != "question header file" {
		t.Error("question.h was not compressed into answer file")
	}
}
