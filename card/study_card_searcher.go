package card

type StudyCardSearcher struct {
	newCards      int
	cardsToRemind int
}

func NewStudyCardSearcher(newCards int, cardsToRemind int) *StudyCardSearcher {
	return &StudyCardSearcher{
		newCards:      newCards,
		cardsToRemind: cardsToRemind,
	}
}
