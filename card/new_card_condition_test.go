package card_test

import (
	"testing"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/pipeline"
)

type cardEmitter struct {
	cards []card.BriefCard
}

func (g *cardEmitter) Run(output chan<- card.BriefCard) error {
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
		WithCardsScheduledToRemind(90).
		Generate()

	briefCards := card.ExtractBriefCards(cards)

	cardRepo := card.NewFakeCardRepository(cards...)

	newCardCollector := pipeline.NewItemCollector[*card.Card]()
	existingCardCollector := pipeline.NewItemCollector[card.BriefCard]()

	cardEmitter := pipeline.NewEmitter[card.BriefCard](p, &cardEmitter{cards: briefCards})
	newCardCondition := pipeline.WithCondition[card.BriefCard](p, cardEmitter, card.NewCardCondition(cardRepo))

	pipeline.WithAcceptor[*card.Card](p, pipeline.OnPositiveDecision(newCardCondition), newCardCollector)
	pipeline.WithAcceptor[card.BriefCard](p, pipeline.OnNegativeDecision(newCardCondition), existingCardCollector)

	err := p.SyncRun()

	if err != nil {
		t.Fatal(err)
	}

	if len(newCardCollector.Items) != 80 {
		t.Fatal("Filter should find 80 new cards")
	}

	if len(existingCardCollector.Items) != 90 {
		t.Fatal("Filter should find 90 existing cards")
	}
}
