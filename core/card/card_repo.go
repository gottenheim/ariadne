package card

type CardRepository interface {
	Get(section string, entry string) (*Card, error)
	Save(card *Card) error
	SaveActivities(card *Card) error
}
