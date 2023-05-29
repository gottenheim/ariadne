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

func TestNewCardFilter(t *testing.T) {
	p := pipeline.New()

	cards := card.NewBatchCardGenerator().
		WithNewCards(80).
		WithLearnedCards(90).
		WithRemindedCards(100).
		Generate()

	keysWithActivities := card.ExtractKeysWithActivities(cards)

	cardRepo := card.NewFakeCardRepository(cards...)

	cardCollector := pipeline.NewItemCollector[*card.Card]()

	generator := pipeline.NewEmitter[*card.KeyWithActivities](p, newCardGenerator(keysWithActivities...))
	filter := pipeline.WithFilter[*card.KeyWithActivities](p, generator, card.NewCardFilter(cardRepo))
	pipeline.WithAcceptor[*card.Card](p, filter, cardCollector)

	err := p.SyncRun()
	if err != nil {
		t.Fatal(err)
	}

	if len(cardCollector.Items) != 80 {
		t.Fatal("Filter should find ten cards")
	}
}
