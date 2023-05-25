package card

import (
	"strconv"
	"strings"
)

type FakeCardRepository struct {
	cards map[string]*Card
}

func NewFakeCardRepository(cards ...*Card) CardRepository {
	repo := &FakeCardRepository{}

	for _, card := range cards {
		cardKey := repo.getCardKey(card)
		repo.cards[cardKey] = card
	}

	return repo
}

func (r *FakeCardRepository) Get(key string) (*Card, error) {
	card, _ := r.cards[key]
	return card, nil
}

func (r *FakeCardRepository) Save(card *Card) error {
	cardKey := r.getCardKey(card)
	r.cards[cardKey] = card
	return nil
}

func (r *FakeCardRepository) getCardKey(card *Card) string {
	return strings.Join(card.sections, "/") + strconv.Itoa(card.orderNum)
}
