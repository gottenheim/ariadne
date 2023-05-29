package card

import "math/rand"

type BatchCardGenerator struct {
	newCards               int
	cardsScheduledToRemind int
}

func NewBatchCardGenerator() *BatchCardGenerator {
	return &BatchCardGenerator{}
}

func (g *BatchCardGenerator) WithNewCards(newCards int) *BatchCardGenerator {
	g.newCards = newCards
	return g
}

func (g *BatchCardGenerator) WithCardsScheduledToRemind(cardsScheduledToRemind int) *BatchCardGenerator {
	g.cardsScheduledToRemind = cardsScheduledToRemind
	return g
}

func (g *BatchCardGenerator) Generate() []*Card {
	var cards []*Card

	newCards, cardsScheduledToRemind := g.newCards, g.cardsScheduledToRemind
	cardsTotal := newCards + cardsScheduledToRemind

	for i := 0; i < cardsTotal; i++ {
		index := rand.Int() % 2

		var activities []GenerateActivity

		if index == 0 && newCards > 0 {
			activities = []GenerateActivity{LearnCard}
			newCards--
		} else if index == 1 && cardsScheduledToRemind > 0 {
			activities = []GenerateActivity{LearnCard | CardExecutedMonthAgo, RemindCard | RemindCardScheduledToToday}
			cardsScheduledToRemind--
		}

		if activities != nil {
			card := NewFakeCard().WithKey(Key(i + 1)).WithActivities(activities...).Build()
			cards = append(cards, card)
		} else {
			i--
		}
	}

	return cards
}
