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

func CollectCards(t *testing.T, cards []*card.Card, config study.DailyCardsConfig) (*study.DailyCards, bool) {
	timeSource := datetime.NewFakeTimeSource()
	cardRepo := card.NewFakeCardRepository(cards...)

	cardEmitter := fakeCardEmitter{
		briefCards: card.ExtractBriefCards(cards),
	}

	dailyCardCollector := study.NewDailyCardsCollector(timeSource, cardRepo, &cardEmitter)

	dailyCardCollector.SetConfig(&config)

	dailyCards, err := dailyCardCollector.Collect()

	if err != nil {
		t.Fatal(err)
	}

	for _, newCard := range dailyCards.NewCards {
		isNewCard, err := card.IsNewCardActivities(newCard.Activities())
		if err != nil {
			t.Fatal(err)
		}

		isCardLearnedToday, err := card.IsCardLearnedToday(timeSource, newCard.Activities())
		if err != nil {
			t.Fatal(err)
		}

		if !isNewCard && !isCardLearnedToday {
			t.Fatal("Card is not new and not learned today but recognized as new")
		}
	}

	for _, scheduledCard := range dailyCards.ScheduledCards {
		isScheduledToRemindToday, err := card.IsCardScheduledToRemindToday(timeSource, scheduledCard.Activities())
		if err != nil {
			t.Fatal(err)
		}

		isRemindedToday, err := card.IsCardRemindedToday(timeSource, scheduledCard.Activities())
		if err != nil {
			t.Fatal(err)
		}

		if !isScheduledToRemindToday && !isRemindedToday {
			t.Fatal("Card is not scheduled to remind today and not reminded today but recognized as scheduled to remind")
		}
	}

	return dailyCards, cardEmitter.cancelled
}

func TestCollectingNewDailyCards(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 100, card.LearnCard)).
		Generate()

	result, cancelled := CollectCards(t, cards, study.DailyCardsConfig{NewCardsCount: 10, ScheduledCardsCount: 0})

	if len(result.NewCards) != 10 {
		t.Errorf("Wrong count of new cards collected. Expected: 10, actual: %d", len(result.NewCards))
	}

	if len(result.ScheduledCards) != 0 {
		t.Errorf("Wrong count of scheduled cards collected. Expected: 0, actual: %d", len(result.ScheduledCards))
	}

	if !cancelled {
		t.Errorf("Cards more than required, process should be cancelled")
	}
}

func TestCollectingNewDailyCardsLearnedToday(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards learned today", 100, card.LearnCard|card.CardExecutedToday, card.RemindCard|card.RemindCardScheduledToTomorrow)).
		Generate()

	result, cancelled := CollectCards(t, cards, study.DailyCardsConfig{NewCardsCount: 10, ScheduledCardsCount: 0})

	if len(result.NewCards) != 10 {
		t.Errorf("Wrong count of new cards collected. Expected: 10, actual: %d", len(result.NewCards))
	}

	if len(result.ScheduledCards) != 0 {
		t.Errorf("Wrong count of scheduled cards collected. Expected: 0, actual: %d", len(result.ScheduledCards))
	}

	if !cancelled {
		t.Errorf("Cards more than required, process should be cancelled")
	}
}

func TestCollectingDailyCardsScheduledToToday(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Cards scheduled to today", 100, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToToday)).
		Generate()

	result, cancelled := CollectCards(t, cards, study.DailyCardsConfig{NewCardsCount: 0, ScheduledCardsCount: 10})

	if len(result.NewCards) != 0 {
		t.Errorf("Wrong count of new cards collected. Expected: 0, actual: %d", len(result.NewCards))
	}

	if len(result.ScheduledCards) != 10 {
		t.Errorf("Wrong count of scheduled cards collected. Expected: 10, actual: %d", len(result.ScheduledCards))
	}

	if !cancelled {
		t.Errorf("Cards more than required, process should be cancelled")
	}
}

func TestCollectingDailyCardsScheduledToYesterday(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Cards scheduled to yesterday", 100, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToYesterday)).
		Generate()

	result, cancelled := CollectCards(t, cards, study.DailyCardsConfig{NewCardsCount: 0, ScheduledCardsCount: 10})

	if len(result.NewCards) != 0 {
		t.Errorf("Wrong count of new cards collected. Expected: 0, actual: %d", len(result.NewCards))
	}

	if len(result.ScheduledCards) != 10 {
		t.Errorf("Wrong count of scheduled cards collected. Expected: 10, actual: %d", len(result.ScheduledCards))
	}

	if !cancelled {
		t.Errorf("Cards more than required, process should be cancelled")
	}
}

func TestCollectingDailyCardsRemindedYesterdayAndRescheduledToFuture(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Cards reminded yesterday and scheduled to future", 100,
			card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToYesterday|card.CardExecutedYesterday, card.RemindCard|card.RemindCardScheduledToMonthAhead)).
		Generate()

	result, cancelled := CollectCards(t, cards, study.DailyCardsConfig{NewCardsCount: 0, ScheduledCardsCount: 10})

	if len(result.NewCards) != 0 {
		t.Errorf("Wrong count of new cards collected. Expected: 0, actual: %d", len(result.NewCards))
	}

	if len(result.ScheduledCards) != 0 {
		t.Errorf("Wrong count of scheduled cards collected. Expected: 0, actual: %d", len(result.ScheduledCards))
	}

	if cancelled {
		t.Errorf("Cards less than required, process should not be cancelled")
	}
}

func TestCollectingNewDailyCardsIfNoNewCardsExist(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Cards scheduled to today", 100, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToToday)).
		Generate()

	result, cancelled := CollectCards(t, cards, study.DailyCardsConfig{NewCardsCount: 10, ScheduledCardsCount: 0})

	if len(result.NewCards) != 0 {
		t.Errorf("Wrong count of new cards collected. Expected: 0, actual: %d", len(result.NewCards))
	}

	if len(result.ScheduledCards) != 0 {
		t.Errorf("Wrong count of scheduled cards collected. Expected: 0, actual: %d", len(result.ScheduledCards))
	}

	if cancelled {
		t.Errorf("Cards less than required, process should not be cancelled")
	}
}

func TestCollectingScheduledDailyCardsIfNoScheduledCardsExist(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 100, card.LearnCard)).
		Generate()

	result, cancelled := CollectCards(t, cards, study.DailyCardsConfig{NewCardsCount: 0, ScheduledCardsCount: 10})

	if len(result.NewCards) != 0 {
		t.Errorf("Wrong count of new cards collected. Expected: 0, actual: %d", len(result.NewCards))
	}

	if len(result.ScheduledCards) != 0 {
		t.Errorf("Wrong count of scheduled cards collected. Expected: 0, actual: %d", len(result.ScheduledCards))
	}

	if cancelled {
		t.Errorf("Cards less than required, process should not be cancelled")
	}
}

func TestCollectingNewAndScheduledDailyCards(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 100, card.LearnCard)).
		WithCards(card.NewCardGenerationSpec("Cards scheduled to today", 100, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToToday)).
		WithCards(card.NewCardGenerationSpec("Cards scheduled to tomorrow", 100, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToTomorrow)).
		Generate()

	result, cancelled := CollectCards(t, cards, study.DailyCardsConfig{NewCardsCount: 10, ScheduledCardsCount: 20})

	if len(result.NewCards) != 10 {
		t.Errorf("Wrong count of new cards collected. Expected: 10, actual: %d", len(result.NewCards))
	}

	if len(result.ScheduledCards) != 20 {
		t.Errorf("Wrong count of scheduled cards collected. Expected: 20, actual: %d", len(result.ScheduledCards))
	}

	if !cancelled {
		t.Errorf("Cards more than required, process should be cancelled")
	}
}

func TestCollectingNewAndScheduledDailyCardsIfNotEnoughRelevantCardsExist(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 10, card.LearnCard)).
		WithCards(card.NewCardGenerationSpec("Cards scheduled to today", 10, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToToday)).
		WithCards(card.NewCardGenerationSpec("Cards scheduled to tomorrow", 100, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToTomorrow)).
		Generate()

	result, cancelled := CollectCards(t, cards, study.DailyCardsConfig{NewCardsCount: 20, ScheduledCardsCount: 20})

	if len(result.NewCards) != 10 {
		t.Errorf("Wrong count of new cards collected. Expected: 10, actual: %d", len(result.NewCards))
	}

	if len(result.ScheduledCards) != 10 {
		t.Errorf("Wrong count of scheduled cards collected. Expected: 20, actual: %d", len(result.ScheduledCards))
	}

	if cancelled {
		t.Errorf("Cards less than required, process should not be cancelled")
	}
}

func TestCollectingDailyCardsLearnedTodayAndScheduledToToday(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards learned today", 100, card.LearnCard|card.CardExecutedToday, card.RemindCard|card.RemindCardScheduledToToday)).
		Generate()

	result, cancelled := CollectCards(t, cards, study.DailyCardsConfig{NewCardsCount: 10, ScheduledCardsCount: 10})

	if len(result.NewCards) != 10 {
		t.Errorf("Wrong count of new cards collected. Expected: 10, actual: %d", len(result.NewCards))
	}

	if len(result.ScheduledCards) != 0 {
		t.Errorf("Wrong count of scheduled cards collected. Expected: 0, actual: %d", len(result.ScheduledCards))
	}

	if cancelled {
		t.Errorf("Cards less than required, process should not be cancelled")
	}
}

func TestCollectingForgottenDailyCards(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Forgotten cards", 100,
			card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToToday|card.CardExecutedToday, card.RemindCard|card.RemindCardScheduledToToday)).
		Generate()

	result, cancelled := CollectCards(t, cards, study.DailyCardsConfig{NewCardsCount: 0, ScheduledCardsCount: 10})

	if len(result.NewCards) != 0 {
		t.Errorf("Wrong count of new cards collected. Expected: 0, actual: %d", len(result.NewCards))
	}

	if len(result.ScheduledCards) != 10 {
		t.Errorf("Wrong count of scheduled cards collected. Expected: 10, actual: %d", len(result.ScheduledCards))
	}

	if !cancelled {
		t.Errorf("Cards more than required, process should be cancelled")
	}
}

func TestCollectingDailyCardsLearnedYesterdayAndScheduledToFuture(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards learned tomorrow", 100, card.LearnCard|card.CardExecutedYesterday, card.RemindCard|card.RemindCardScheduledToMonthAhead)).
		Generate()

	result, cancelled := CollectCards(t, cards, study.DailyCardsConfig{NewCardsCount: 10, ScheduledCardsCount: 10})

	if len(result.NewCards) != 0 {
		t.Errorf("Wrong count of new cards collected. Expected: 0, actual: %d", len(result.NewCards))
	}

	if len(result.ScheduledCards) != 0 {
		t.Errorf("Wrong count of scheduled cards collected. Expected: 0, actual: %d", len(result.ScheduledCards))
	}

	if cancelled {
		t.Errorf("Cards less than required, process should not be cancelled")
	}
}

func TestCollectingDailyCardsScheduledToYesterdayAndRemindedYesterday(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Cards scheduled to yesterday and reminded yesterday", 100,
			card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToYesterday|card.CardExecutedYesterday)).
		Generate()

	result, cancelled := CollectCards(t, cards, study.DailyCardsConfig{NewCardsCount: 10, ScheduledCardsCount: 10})

	if len(result.NewCards) != 0 {
		t.Errorf("Wrong count of new cards collected. Expected: 0, actual: %d", len(result.NewCards))
	}

	if len(result.ScheduledCards) != 0 {
		t.Errorf("Wrong count of scheduled cards collected. Expected: 0, actual: %d", len(result.ScheduledCards))
	}

	if cancelled {
		t.Errorf("Cards less than required, process should not be cancelled")
	}
}
