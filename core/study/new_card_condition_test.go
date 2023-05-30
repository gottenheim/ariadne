package study_test

import (
	"context"
	"testing"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/core/study"
	"github.com/gottenheim/ariadne/libraries/datetime"
	"github.com/gottenheim/ariadne/libraries/pipeline"
)

type cardEmitter struct {
	cards []card.BriefCard
}

func (g *cardEmitter) Run(ctx context.Context, output chan<- card.BriefCard) error {
	for _, card := range g.cards {
		output <- card
	}

	return nil
}

func TestNewCardFilter(t *testing.T) {
	p := pipeline.New()

	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 80, card.LearnCard)).
		WithCards(card.NewCardGenerationSpec("Cards learned today", 90, card.LearnCard|card.CardExecutedToday, card.RemindCard|card.RemindCardScheduledToTomorrow)).
		WithCards(card.NewCardGenerationSpec("Cards learned yesterday", 100, card.LearnCard|card.CardExecutedYesterday, card.RemindCard|card.RemindCardScheduledToTomorrow)).
		Generate()

	briefCards := card.ExtractBriefCards(cards)

	timeSource := datetime.NewFakeTimeSource()
	cardRepo := card.NewFakeCardRepository(cards...)

	newCardCollector := pipeline.NewItemCollector[*card.Card]()
	existingCardCollector := pipeline.NewItemCollector[card.BriefCard]()

	cardEmitter := pipeline.NewEmitter[card.BriefCard](p, &cardEmitter{cards: briefCards})
	newCardCondition := pipeline.WithCondition[card.BriefCard](p, cardEmitter, study.NewCardCondition(timeSource, cardRepo))

	pipeline.WithAcceptor[*card.Card](p, pipeline.OnPositiveDecision(newCardCondition), newCardCollector)
	pipeline.WithAcceptor[card.BriefCard](p, pipeline.OnNegativeDecision(newCardCondition), existingCardCollector)

	err := p.SyncRun()

	if err != nil {
		t.Fatal(err)
	}

	if len(newCardCollector.Items) != 170 {
		t.Fatal("Filter should find 170 new cards and cards learned today")
	}

	if len(existingCardCollector.Items) != 100 {
		t.Fatal("Filter should find 100 cards learned earlier")
	}
}