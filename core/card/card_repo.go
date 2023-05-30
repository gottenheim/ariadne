package card

type CardRepository interface {
	Get(key Key) (*Card, error)
	Save(card *Card) error
}
