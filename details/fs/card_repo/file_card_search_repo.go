package card_repo

import "github.com/gottenheim/ariadne/card"

func (r *FileCardRepository) FindCardsForStudy(newCards int, cardsToRemind int) ([]*card.Card, error) {
	// var foundCards []*Card

	// err := afero.Walk(r.fs, r.baseDir, func(filePath string, info os.FileInfo, err error) error {
	// 	if info == nil || !info.IsDir() {
	// 		return nil
	// 	}

	// 	isCardDir, err := r.isCardDir(filePath)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if !isCardDir {
	// 		return nil
	// 	}

	// 	cardDir := filePath

	// 	activities, err := r.ReadCardActivities(cardDir)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	newCardsLeft := newCards
	// 	cardsToRemindLeft := cardsToRemind

	// 	loadCard := false

	// 	isNewCard, err := IsNewCard(activities)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if isNewCard {
	// 		newCardsLeft--
	// 		loadCard = true
	// 	} else {

	// 	}

	// 	if newCardsLeft == 0 && cardsToRemindLeft == 0 {
	// 		return io.EOF
	// 	}

	// 	return filepath.SkipDir
	// })

	// if err != nil && err != io.EOF {
	// 	return nil, err
	// }

	return nil, nil
}
