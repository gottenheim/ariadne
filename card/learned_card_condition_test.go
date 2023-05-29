package card_test

import (
	"testing"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/datetime"
	"github.com/gottenheim/ariadne/details/pipeline"
)

func TestCardLearnedTodayFilter(t *testing.T) {
	p := pipeline.New()

	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Cards learned today", 80, card.LearnCard|card.CardExecutedToday, card.RemindCard|card.RemindCardScheduledToTomorrow)).
		WithCards(card.NewCardGenerationSpec("Cards not learned today", 90, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToTomorrow)).
		Generate()

	briefCards := card.ExtractBriefCards(cards)

	timeSource := datetime.NewFakeTimeSource()
	cardRepo := card.NewFakeCardRepository(cards...)

	learnedCardCollector := pipeline.NewItemCollector[*card.Card]()
	notLearnedCardCollector := pipeline.NewItemCollector[card.BriefCard]()

	cardEmitter := pipeline.NewEmitter[card.BriefCard](p, &cardEmitter{cards: briefCards})
	learnedCardCondition := pipeline.WithCondition[card.BriefCard](p, cardEmitter, card.LearnedCardCondition(timeSource, cardRepo))

	pipeline.WithAcceptor[*card.Card](p, pipeline.OnPositiveDecision(learnedCardCondition), learnedCardCollector)
	pipeline.WithAcceptor[card.BriefCard](p, pipeline.OnNegativeDecision(learnedCardCondition), notLearnedCardCollector)

	err := p.SyncRun()

	if err != nil {
		t.Fatal(err)
	}

	if len(learnedCardCollector.Items) != 80 {
		t.Fatal("Filter should find 80 cards learned today")
	}

	if len(notLearnedCardCollector.Items) != 90 {
		t.Fatal("Filter should find 90 cards not learned today")
	}
}
