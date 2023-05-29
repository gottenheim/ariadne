package card_test

import (
	"testing"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/pipeline"
)

type cardGenerator struct {
	cards []*card.KeyWithActivities
}

func newCardGenerator(cards ...*card.KeyWithActivities) *cardGenerator {
	return &cardGenerator{
		cards: cards,
	}
}

func (g *cardGenerator) Run(output chan<- *card.KeyWithActivities) error {
	defer func() {
		close(output)
	}()

	for _, card := range g.cards {
		output <- card
	}

	return nil
}

type cardFilteringResult struct {
	cards []*card.Card
}

type cardAccumulator struct {
	result *cardFilteringResult
}

func newCardAccumulator(result *cardFilteringResult) *cardAccumulator {
	return &cardAccumulator{
		result: result,
	}
}

func (g *cardAccumulator) Run(input <-chan *card.Card, output chan<- interface{}) error {
	for {
		card, ok := <-input
		if !ok {
			break
		}
		g.result.cards = append(g.result.cards, card)
	}

	return nil
}

func TestNewCardFilter(t *testing.T) {
	p := pipeline.New()

	cards := card.NewBatchCardGenerator().
		WithNewCards(80).
		WithLearnedCards(90).
		WithRemindedCards(100).
		Generate()

	keysWithActivities := card.ExtractKeysWithActivities(cards)

	cardRepo := card.NewFakeCardRepository(cards...)

	filteringResult := &cardFilteringResult{}

	generator := pipeline.NewGenerator[*card.KeyWithActivities](p, newCardGenerator(keysWithActivities...))
	filter := pipeline.WithFilter[*card.KeyWithActivities](p, generator, card.NewCardFilter(cardRepo))
	pipeline.WithFilter[*card.Card, interface{}](p, filter, newCardAccumulator(filteringResult))

	err := p.SyncRun()
	if err != nil {
		t.Fatal(err)
	}

	if len(filteringResult.cards) != 80 {
		t.Fatal("Filter should find ten cards")
	}
}
