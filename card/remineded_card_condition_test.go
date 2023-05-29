package card_test

import (
	"testing"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/datetime"
	"github.com/gottenheim/ariadne/details/pipeline"
)

func TestCardRemindedTodayFilter(t *testing.T) {
	p := pipeline.New()

	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Cards reminded today", 80, card.LearnCard|card.CardExecutedMonthAgo,
			card.RemindCard|card.RemindCardScheduledToToday|card.CardExecutedToday, card.RemindCard|card.RemindCardScheduledToMonthAhead)).
		WithCards(card.NewCardGenerationSpec("Cards not reminded to today", 90, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToTomorrow)).
		Generate()

	briefCards := card.ExtractBriefCards(cards)

	timeSource := datetime.NewFakeTimeSource()
	cardRepo := card.NewFakeCardRepository(cards...)

	remindedCardCollector := pipeline.NewItemCollector[*card.Card]()
	notRemindedCardCollector := pipeline.NewItemCollector[card.BriefCard]()

	cardEmitter := pipeline.NewEmitter[card.BriefCard](p, &cardEmitter{cards: briefCards})
	newCardCondition := pipeline.WithCondition[card.BriefCard](p, cardEmitter, card.RemindedCardCondition(timeSource, cardRepo))

	pipeline.WithAcceptor[*card.Card](p, pipeline.OnPositiveDecision(newCardCondition), remindedCardCollector)
	pipeline.WithAcceptor[card.BriefCard](p, pipeline.OnNegativeDecision(newCardCondition), notRemindedCardCollector)

	err := p.SyncRun()

	if err != nil {
		t.Fatal(err)
	}

	if len(remindedCardCollector.Items) != 80 {
		t.Fatal("Filter should find 80 cards reminded today")
	}

	if len(notRemindedCardCollector.Items) != 90 {
		t.Fatal("Filter should find 90 cards not reminded today")
	}
}
