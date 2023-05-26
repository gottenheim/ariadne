package card

type FakeCardRepository struct {
	cards map[Key]*Card
}

func NewFakeCardRepository(cards ...*Card) CardRepository {
	repo := &FakeCardRepository{
		cards: map[Key]*Card{},
	}

	for _, card := range cards {
		repo.cards[card.Key()] = card
	}

	return repo
}

func (r *FakeCardRepository) Get(key Key) (*Card, error) {
	card, _ := r.cards[key]
	return card, nil
}

func (r *FakeCardRepository) Save(card *Card) error {
	r.cards[card.Key()] = card
	return nil
}
