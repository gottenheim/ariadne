package card

type CardRepository interface {
	Get(key string) (*Card, error)
	Save(card *Card) error
}
