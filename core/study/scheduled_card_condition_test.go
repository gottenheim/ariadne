package study_test

import (
	"testing"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/core/study"
	"github.com/gottenheim/ariadne/libraries/datetime"
	"github.com/gottenheim/ariadne/libraries/pipeline"
)

func TestCardScheduledToTodayCondition_ShouldCountCardsScheduledToToday(t *testing.T) {
	p := pipeline.New()

	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Cards scheduled to today", 80, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToToday)).
		WithCards(card.NewCardGenerationSpec("Cards scheduled to future", 90, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToTomorrow)).
		Generate()

	briefCards := card.ExtractBriefCards(cards)

	timeSource := datetime.NewFakeTimeSource()
	cardRepo := card.NewFakeCardRepository(cards...)

	scheduledCardCollector := pipeline.NewItemCollector[*card.Card]()
	notScheduledCardCollector := pipeline.NewItemCollector[card.BriefCard]()

	cardEmitter := pipeline.NewEmitter[card.BriefCard](p, &cardEmitter{cards: briefCards})
	newCardCondition := pipeline.WithCondition[card.BriefCard](p, cardEmitter, study.ScheduledCardCondition(timeSource, cardRepo))

	pipeline.WithAcceptor[*card.Card](p, pipeline.OnPositiveDecision(newCardCondition), scheduledCardCollector)
	pipeline.WithAcceptor[card.BriefCard](p, pipeline.OnNegativeDecision(newCardCondition), notScheduledCardCollector)

	err := p.SyncRun()

	if err != nil {
		t.Fatal(err)
	}

	if len(scheduledCardCollector.Items) != 80 {
		t.Fatal("Filter should find 80 cards scheduled to today")
	}

	if len(notScheduledCardCollector.Items) != 90 {
		t.Fatal("Filter should find 90 cards scheduled to future")
	}
}

func TestCardScheduledToTodayCondition_ShouldCountCardsScheduledAndRemindedToday(t *testing.T) {
	p := pipeline.New()

	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Cards scheduled and reminded today", 80, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToToday|card.CardExecutedToday)).
		WithCards(card.NewCardGenerationSpec("Cards scheduled to future", 90, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToTomorrow)).
		Generate()

	briefCards := card.ExtractBriefCards(cards)

	timeSource := datetime.NewFakeTimeSource()
	cardRepo := card.NewFakeCardRepository(cards...)

	scheduledCardCollector := pipeline.NewItemCollector[*card.Card]()
	notScheduledCardCollector := pipeline.NewItemCollector[card.BriefCard]()

	cardEmitter := pipeline.NewEmitter[card.BriefCard](p, &cardEmitter{cards: briefCards})
	newCardCondition := pipeline.WithCondition[card.BriefCard](p, cardEmitter, study.ScheduledCardCondition(timeSource, cardRepo))

	pipeline.WithAcceptor[*card.Card](p, pipeline.OnPositiveDecision(newCardCondition), scheduledCardCollector)
	pipeline.WithAcceptor[card.BriefCard](p, pipeline.OnNegativeDecision(newCardCondition), notScheduledCardCollector)

	err := p.SyncRun()

	if err != nil {
		t.Fatal(err)
	}

	if len(scheduledCardCollector.Items) != 80 {
		t.Fatal("Filter should find 80 cards scheduled and reminded today")
	}

	if len(notScheduledCardCollector.Items) != 90 {
		t.Fatal("Filter should find 90 cards scheduled to future")
	}
}

func TestCardScheduledToTodayCondition_ShouldCountCardsScheduledToYesterdayAndRemindedToday(t *testing.T) {
	p := pipeline.New()

	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Cards scheduled to yesterday and reminded today", 80, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToYesterday|card.CardExecutedToday)).
		WithCards(card.NewCardGenerationSpec("Cards scheduled to future", 90, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToTomorrow)).
		Generate()

	briefCards := card.ExtractBriefCards(cards)

	timeSource := datetime.NewFakeTimeSource()
	cardRepo := card.NewFakeCardRepository(cards...)

	scheduledCardCollector := pipeline.NewItemCollector[*card.Card]()
	notScheduledCardCollector := pipeline.NewItemCollector[card.BriefCard]()

	cardEmitter := pipeline.NewEmitter[card.BriefCard](p, &cardEmitter{cards: briefCards})
	newCardCondition := pipeline.WithCondition[card.BriefCard](p, cardEmitter, study.ScheduledCardCondition(timeSource, cardRepo))

	pipeline.WithAcceptor[*card.Card](p, pipeline.OnPositiveDecision(newCardCondition), scheduledCardCollector)
	pipeline.WithAcceptor[card.BriefCard](p, pipeline.OnNegativeDecision(newCardCondition), notScheduledCardCollector)

	err := p.SyncRun()

	if err != nil {
		t.Fatal(err)
	}

	if len(scheduledCardCollector.Items) != 80 {
		t.Fatal("Filter should find 80 cards scheduled to yesterday and reminded today")
	}

	if len(notScheduledCardCollector.Items) != 90 {
		t.Fatal("Filter should find 90 cards scheduled to future")
	}
}

func TestCardScheduledToTodayCondition_ShouldCountCardsScheduledToYesterday(t *testing.T) {
	p := pipeline.New()

	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Cards scheduled to yesterday", 80, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToYesterday)).
		WithCards(card.NewCardGenerationSpec("Cards scheduled to future", 90, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToTomorrow)).
		Generate()

	briefCards := card.ExtractBriefCards(cards)

	timeSource := datetime.NewFakeTimeSource()
	cardRepo := card.NewFakeCardRepository(cards...)

	scheduledCardCollector := pipeline.NewItemCollector[*card.Card]()
	notScheduledCardCollector := pipeline.NewItemCollector[card.BriefCard]()

	cardEmitter := pipeline.NewEmitter[card.BriefCard](p, &cardEmitter{cards: briefCards})
	newCardCondition := pipeline.WithCondition[card.BriefCard](p, cardEmitter, study.ScheduledCardCondition(timeSource, cardRepo))

	pipeline.WithAcceptor[*card.Card](p, pipeline.OnPositiveDecision(newCardCondition), scheduledCardCollector)
	pipeline.WithAcceptor[card.BriefCard](p, pipeline.OnNegativeDecision(newCardCondition), notScheduledCardCollector)

	err := p.SyncRun()

	if err != nil {
		t.Fatal(err)
	}

	if len(scheduledCardCollector.Items) != 80 {
		t.Fatal("Filter should find 80 cards scheduled to yesterday")
	}

	if len(notScheduledCardCollector.Items) != 90 {
		t.Fatal("Filter should find 90 cards scheduled to future")
	}
}
