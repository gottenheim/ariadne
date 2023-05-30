package study_test

import (
	"context"
	"testing"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/core/study"
	"github.com/gottenheim/ariadne/libraries/datetime"
)

type fakeCardEmitter struct {
	briefCards []card.BriefCard
	cancelled  bool
}

func (e *fakeCardEmitter) Run(ctx context.Context, output chan<- card.BriefCard) error {
	e.cancelled = false
	for _, card := range e.briefCards {
		select {
		case <-ctx.Done():
			e.cancelled = true
			break
		case output <- card:
		}

	}
	return nil
}

func TestCollectingDailyCards(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 100, card.LearnCard)).
		WithCards(card.NewCardGenerationSpec("Scheduled cards", 200, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToToday)).
		Generate()

	timeSource := datetime.NewFakeTimeSource()
	cardRepo := card.NewFakeCardRepository(cards...)

	cardEmitter := fakeCardEmitter{
		briefCards: card.ExtractBriefCards(cards),
	}

	dailyCardCollector := study.NewDailyCardsCollector(timeSource, cardRepo, &cardEmitter)

	config := &study.DailyCardsConfig{
		NewCardsCount:       10,
		ScheduledCardsCount: 20,
	}

	dailyCardCollector.SetConfig(config)

	dailyCards, err := dailyCardCollector.Collect()

	if err != nil {
		t.Fatal(err)
	}

	if len(dailyCards.NewCards) != config.NewCardsCount {
		t.Errorf("Wrong count of new cards collected. Expected: %d, actual: %d", config.NewCardsCount, len(dailyCards.NewCards))
	}

	if len(dailyCards.ScheduledCards) != config.ScheduledCardsCount {
		t.Errorf("Wrong count of scheduled cards collected. Expected: %d, actual: %d", config.ScheduledCardsCount, len(dailyCards.ScheduledCards))
	}

	if !cardEmitter.cancelled {
		t.Error("Card emitter should be cancelled")
	}
}
