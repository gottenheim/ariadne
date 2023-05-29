package card

import "math/rand"

type CardGenerationSpec struct {
	name       string
	count      int
	activities []GenerateActivity
}

func NewCardGenerationSpec(name string, count int, activities ...GenerateActivity) CardGenerationSpec {
	return CardGenerationSpec{
		name:       name,
		count:      count,
		activities: activities,
	}
}

type BatchCardGenerator struct {
	specs []CardGenerationSpec
}

func NewBatchCardGenerator() *BatchCardGenerator {
	return &BatchCardGenerator{}
}

func (g *BatchCardGenerator) WithCards(spec CardGenerationSpec) *BatchCardGenerator {
	g.specs = append(g.specs, spec)
	return g
}

func (g *BatchCardGenerator) Generate() []*Card {
	var cards []*Card

	cardsTotal := g.getCardsTotal()

	for i := 0; i < cardsTotal; i++ {
		index := rand.Int() % len(g.specs)

		spec := &g.specs[index]

		if spec.count > 0 {
			card := NewFakeCard().WithKey(Key(i + 1)).WithActivities(spec.activities...).Build()
			cards = append(cards, card)
			spec.count--
		} else {
			i--
		}
	}

	return cards
}

func (g *BatchCardGenerator) getCardsTotal() int {
	cardsTotal := 0
	for _, spec := range g.specs {
		cardsTotal += spec.count
	}
	return cardsTotal
}
