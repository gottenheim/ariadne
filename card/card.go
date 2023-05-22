package card

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/gottenheim/ariadne/archive"
	"github.com/spf13/afero"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Card struct {
	fs        afero.Fs
	ioStreams genericclioptions.IOStreams
	config    *Config
}

func New(fs afero.Fs, config *Config, ioStreams genericclioptions.IOStreams) *Card {
	return &Card{
		fs:        fs,
		ioStreams: ioStreams,
		config:    config,
	}
}

/*
	 	Preconditions:
	 	- cards directory exists and writable
		- template directory exists and readable
	 	Postconditions:
		- new card directory created in cards directory
		- template files copied to that directory
*/
func (c *Card) CreateCard(cardsDirPath string, templateDirPath string) (string, error) {
	cardDirPath, err := c.getNextCardDirPath(cardsDirPath)

	if err != nil {
		return "", err
	}

	err = c.createCardDirectory(cardDirPath)

	if err != nil {
		return "", err
	}

	err = c.copyTemplateFilesToCardDirectory(cardDirPath, templateDirPath)

	if err != nil {
		return "", err
	}

	return cardDirPath, nil
}

/*
	 	Preconditions:
	 	- card directory exists and writable
	 	Postconditions:
		- all code artifacts (answers) compressed and saved to archive file
*/
func (c *Card) PackAnswer(cardDirPath string) error {
	err := c.removeAnswerFile(cardDirPath)
	if err != nil {
		return err
	}

	return c.putCardFilesIntoArchive(cardDirPath)
}

/*
	 	Preconditions:
	 	- card directory exists and writable
		- answer archive file exists
	 	Postconditions:
		- all code artifacts (answers) extracted and saved to card directory
*/
func (c *Card) UnpackAnswer(cardDirPath string) error {
	return c.extractCardFilesToCardDirectory(cardDirPath)
}

/*
	 	Preconditions:
	 	- card directory exists and readable
		- answer archive file exists
	 	Postconditions:
		- all code artifacts (answers) extracted and sent to stdout
*/
func (c *Card) ShowAnswer(cardDirPath string) error {
	return c.extractAndDisplayCardFiles(cardDirPath)
}

func (c *Card) getNextCardDirPath(cardsDirPath string) (string, error) {
	maxCardNumber := 0

	err := afero.Walk(c.fs, cardsDirPath, func(filePath string, info os.FileInfo, err error) error {
		isDir, _ := afero.IsDir(c.fs, filePath)

		if isDir {
			cardDir := path.Base(filePath)
			cardNumber, err := strconv.Atoi(cardDir)
			if err != nil {
				return nil
			}

			if cardNumber > maxCardNumber {
				maxCardNumber = cardNumber
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	cardDirPath := path.Join(cardsDirPath, fmt.Sprintf("%d", maxCardNumber+1))

	return cardDirPath, nil
}

func (c *Card) createCardDirectory(cardDirPath string) error {
	return c.fs.MkdirAll(cardDirPath, os.ModePerm)
}

func (c *Card) copyTemplateFilesToCardDirectory(cardDirPath string, templateDirPath string) error {
	return afero.Walk(c.fs, templateDirPath, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			srcFileContents, err := afero.ReadFile(c.fs, filePath)
			if err != nil {
				return err
			}
			fileName := path.Base(filePath)
			dstFilePath := path.Join(cardDirPath, fileName)
			afero.WriteFile(c.fs, dstFilePath, srcFileContents, os.ModePerm)
		}
		return nil
	})
}

func (c *Card) getAnswerFilePath(cardDirPath string) string {
	return path.Join(cardDirPath, c.config.AnswerFileName)
}

func (c *Card) removeAnswerFile(cardDirPath string) error {
	answerFilePath := c.getAnswerFilePath(cardDirPath)
	exists, err := afero.Exists(c.fs, answerFilePath)
	if err != nil {
		return err
	}
	if exists {
		return c.fs.Remove(answerFilePath)
	}
	return nil
}

func (c *Card) putCardFilesIntoArchive(cardDirPath string) error {
	archiveWriter := archive.NewWriter()
	err := archiveWriter.AddDir(c.fs, cardDirPath)
	if err != nil {
		return err
	}

	answerFilePath := c.getAnswerFilePath(cardDirPath)

	buf, err := archiveWriter.Buffer()
	if err != nil {
		return err
	}

	return afero.WriteFile(c.fs, answerFilePath, buf.Bytes(), os.ModePerm)
}

func (c *Card) extractCardFilesToCardDirectory(cardDirPath string) error {
	answerFilePath := c.getAnswerFilePath(cardDirPath)

	answerFileContents, err := afero.ReadFile(c.fs, answerFilePath)

	if err != nil {
		return err
	}

	return archive.Uncompress(bytes.NewReader(answerFileContents), c.fs, cardDirPath)
}

func (c *Card) extractAndDisplayCardFiles(cardDirPath string) error {
	answerFilePath := c.getAnswerFilePath(cardDirPath)

	answerFileContents, err := afero.ReadFile(c.fs, answerFilePath)

	if err != nil {
		return err
	}

	files, err := archive.GetFiles(bytes.NewReader(answerFileContents))

	if err != nil {
		return err
	}

	for fileName, fileContents := range files {
		fmt.Fprintf(c.ioStreams.Out, "---- %s ----\n", fileName)
		fmt.Fprintf(c.ioStreams.Out, "%s\n", string(fileContents))
	}

	return nil
}
