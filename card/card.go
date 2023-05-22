package card

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/spf13/afero"
)

type Card struct {
	fs     afero.Fs
	config *Config
}

func New(fs afero.Fs, config *Config) *Card {
	return &Card{
		fs:     fs,
		config: config,
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
func (c *Card) CreateNew(cardsDirPath string) (string, error) {
	cardDirPath, err := c.getNextCardDirPath(cardsDirPath)

	if err != nil {
		return "", err
	}

	err = c.createCardDirectory(cardDirPath)

	if err != nil {
		return "", err
	}

	err = c.copyTemplateFilesToCardDirectory(cardDirPath)

	if err != nil {
		return "", err
	}

	return cardDirPath, nil
}

func (c *Card) getNextCardDirPath(cardsDirPath string) (string, error) {
	maxCardNumber := 0

	err := afero.Walk(c.fs, cardsDirPath, func(filePath string, info os.FileInfo, err error) error {
		if path.Base(filePath) == c.config.AnswerFileName {
			cardDir := path.Base(path.Dir(filePath))
			cardNumber, err := strconv.Atoi(cardDir)
			if err != nil {
				return err
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

func (c *Card) copyTemplateFilesToCardDirectory(cardDirPath string) error {
	return afero.Walk(c.fs, c.config.TemplateDir, func(filePath string, info os.FileInfo, err error) error {
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
