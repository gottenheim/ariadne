package card

type UserInteractor interface {
	ShowAnswer(card *Card) error
}
