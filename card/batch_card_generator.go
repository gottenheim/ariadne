package card

import "math/rand"

type BatchCardGenerator struct {
	newCards      int
	learnedCards  int
	remindedCards int
}

func NewBatchCardGenerator() *BatchCardGenerator {
	return &BatchCardGenerator{}
}

func (g *BatchCardGenerator) WithNewCards(newCards int) *BatchCardGenerator {
	g.newCards = newCards
	return g
}

func (g *BatchCardGenerator) WithLearnedCards(learnedCards int) *BatchCardGenerator {
	g.learnedCards = learnedCards
	return g
}

func (g *BatchCardGenerator) WithRemindedCards(remindedCards int) *BatchCardGenerator {
	g.remindedCards = remindedCards
	return g
}

func (g *BatchCardGenerator) Generate() []*Card {
	var cards []*Card

	newCards, learnedCards, remindedCards := g.newCards, g.learnedCards, g.remindedCards
	cardsTotal := newCards + learnedCards + remindedCards

	for i := 0; i < cardsTotal; i++ {
		index := rand.Int() % 3

		var activities []GenerateActivity

		if index == 0 && newCards > 0 {
			activities = []GenerateActivity{LearnCard}
			newCards--
		} else if index == 1 && learnedCards > 0 {
			activities = []GenerateActivity{LearnCard | CardExecutedYesterday, RemindCard | RemindCardScheduledToTomorrow}
			learnedCards--
		} else if index == 2 && remindedCards > 0 {
			activities = []GenerateActivity{LearnCard | CardExecutedMonthAgo, RemindCard | RemindCardScheduledToYesterday | CardExecutedToday, RemindCard | RemindCardScheduledToMonthAhead}
			remindedCards--
		}

		if activities != nil {
			card := NewFakeCard().WithKey(Key(i + 1)).WithActivities(activities...).Build()
			cards = append(cards, card)
		}
	}

	return cards
}
