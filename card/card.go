package card

import (
	"bytes"
	"errors"

	"github.com/gottenheim/ariadne/archive"
)

const AnswerArtifactName = "answer.tgz"

type Card struct {
	sections  []string
	orderNum  int
	artifacts []CardArtifact
	progress  *CardProgress
}

func NewCard(sections []string, orderNum int, artifacts []CardArtifact) *Card {
	return &Card{
		sections:  sections,
		orderNum:  orderNum,
		artifacts: artifacts,
		progress: &CardProgress{
			status: New,
		},
	}
}

func (c *Card) Sections() []string {
	return c.sections
}

func (c *Card) OrderNumber() int {
	return c.orderNum
}

func (c *Card) SetOrderNumber(orderNum int) {
	c.orderNum = orderNum
}

func (c *Card) HasOrderNumber() bool {
	return c.orderNum > 0
}

func (c *Card) Artifacts() []CardArtifact {
	return c.artifacts
}

func (c *Card) SetProgress(progress *CardProgress) {
	c.progress = progress
}

func (c *Card) Progress() *CardProgress {
	return c.progress
}

func (c *Card) FindArtifactByName(name string) *CardArtifact {
	for _, artifact := range c.artifacts {
		if artifact.Name() == name {
			return &artifact
		}
	}
	return nil
}

func (c *Card) FindAnswerArtifact() *CardArtifact {
	return c.FindArtifactByName(AnswerArtifactName)
}

func (c *Card) CompressAnswer() error {
	c.removeAnswerArtifact()

	compressedAnswer, err := c.compressArtifacts()
	if err != nil {
		return err
	}

	c.addAnswerArtifact(compressedAnswer)

	return nil
}

func (c *Card) ExtractAnswer() error {
	files, err := c.getAnswerFiles()
	if err != nil {
		return err
	}

	c.removeAllArtifactsExceptAnswer()

	for name, content := range files {
		c.addArtifact(name, content)
	}

	return nil
}

func (c *Card) removeAnswerArtifact() {
	c.removeArtifact(AnswerArtifactName)
}

func (c *Card) removeArtifact(artifactName string) {
	for index, artifact := range c.artifacts {
		if artifact.Name() == artifactName {
			c.artifacts = append(c.artifacts[:index], c.artifacts[index+1:]...)
			break
		}
	}
}

func (c *Card) removeAllArtifactsExceptAnswer() {
	var newArtifacts []CardArtifact

	answerArtifact := c.FindAnswerArtifact()
	if answerArtifact != nil {
		newArtifacts = append(newArtifacts, *answerArtifact)
	}

	c.artifacts = newArtifacts
}

func (c *Card) compressArtifacts() ([]byte, error) {
	artifacts := c.getArtifactsAsMap()

	writer := archive.NewWriter()
	writer.AddFiles(artifacts)
	buf, err := writer.Buffer()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *Card) getArtifactsAsMap() map[string][]byte {
	artifacts := map[string][]byte{}

	for _, artifact := range c.artifacts {
		artifacts[artifact.name] = artifact.content
	}

	return artifacts
}

func (c *Card) addAnswerArtifact(contents []byte) {
	c.addArtifact(AnswerArtifactName, contents)
}

func (c *Card) addArtifact(name string, contents []byte) {
	c.artifacts = append(c.artifacts, NewCardArtifact(name, contents))
}

func (c *Card) getAnswerFiles() (map[string][]byte, error) {
	answerArtifact := c.FindAnswerArtifact()

	if answerArtifact == nil {
		return nil, errors.New("Card doesn't contain answer artifact")
	}

	files, err := archive.GetFiles(bytes.NewReader(answerArtifact.content))

	if err != nil {
		return nil, err
	}

	return files, nil
}
