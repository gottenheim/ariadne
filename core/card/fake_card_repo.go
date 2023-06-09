package card

import (
	"fmt"
)

type FakeCardRepository struct {
	cards map[string]*Card
}

func NewFakeCardRepository(cards ...*Card) CardRepository {
	repo := &FakeCardRepository{
		cards: map[string]*Card{},
	}

	for _, card := range cards {
		cardKey := repo.getFullCardKey(card.Section(), card.Entry())
		repo.cards[cardKey] = card
	}

	return repo
}

func (r *FakeCardRepository) Get(section string, entry string) (*Card, error) {
	key := r.getFullCardKey(section, entry)
	card, ok := r.cards[key]
	if !ok {
		return nil, fmt.Errorf("Card with key %s could not be found", key)
	}
	return FromExisting(card.section, card.entry, card.artifacts, card.activities), nil
}

func (r *FakeCardRepository) Save(card *Card) error {
	key := r.getFullCardKey(card.Section(), card.Entry())
	r.cards[key] = card
	return nil
}

func (r *FakeCardRepository) SaveActivities(card *Card) error {
	key := r.getFullCardKey(card.Section(), card.Entry())
	existingCard, ok := r.cards[key]
	if !ok {
		return fmt.Errorf("Card with key %s could not be found", key)
	}
	existingCard.SetActivities(card.Activities())
	return nil
}

func (r *FakeCardRepository) getFullCardKey(section string, entry string) string {
	return fmt.Sprintf("%s/%s", section, entry)
}
