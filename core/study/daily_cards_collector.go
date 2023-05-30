package study

import (
	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/libraries/datetime"
	"github.com/gottenheim/ariadne/libraries/pipeline"
)

type DailyCardsConfig struct {
	NewCardsCount       int
	ScheduledCardsCount int
}

type DailyCards struct {
	NewCards       []*card.Card
	ScheduledCards []*card.Card
}

type DailyCardsCollector struct {
	timeSource  datetime.TimeSource
	cardRepo    card.CardRepository
	cardEmitter pipeline.Emitter[card.BriefCard]
	config      *DailyCardsConfig
}

func NewDailyCardsCollector(timeSource datetime.TimeSource, cardRepo card.CardRepository, cardEmitter pipeline.Emitter[card.BriefCard]) *DailyCardsCollector {
	return &DailyCardsCollector{
		timeSource:  timeSource,
		cardRepo:    cardRepo,
		cardEmitter: cardEmitter,
	}
}

func (c *DailyCardsCollector) SetConfig(config *DailyCardsConfig) {
	c.config = config
}

func (c *DailyCardsCollector) Collect() (*DailyCards, error) {
	newCardsCollector := pipeline.NewPassingItemCollector[*card.Card]()
	scheduledCardsCollector := pipeline.NewPassingItemCollector[*card.Card]()

	p := pipeline.New()
	cardEmissionStep := pipeline.NewEmitter(p, c.cardEmitter)
	isNewCardStep := pipeline.WithCondition[card.BriefCard](p, cardEmissionStep, NewCardCondition(c.timeSource, c.cardRepo))
	limitNewCardsStep := pipeline.WithFilter(p, pipeline.OnPositiveDecision(isNewCardStep), pipeline.Limit[*card.Card](c.config.NewCardsCount))
	pipeline.WithAcceptor(p, pipeline.OnNegativeDecision(isNewCardStep), pipeline.DevNull[card.BriefCard]())
	collectNewCardsStep := pipeline.WithFilter[*card.Card, *card.Card](p, limitNewCardsStep, newCardsCollector)
	countNewCardsStep := pipeline.WithFilter[*card.Card, int](p, collectNewCardsStep, pipeline.NewCounter[*card.Card]())

	isScheduledCardStep := pipeline.WithCondition(p, pipeline.OnNegativeDecision(isNewCardStep), ScheduledCardCondition(c.timeSource, c.cardRepo))
	limitScheduledCardsStep := pipeline.WithFilter(p, pipeline.OnPositiveDecision(isScheduledCardStep), pipeline.Limit[*card.Card](c.config.ScheduledCardsCount))
	pipeline.WithAcceptor(p, pipeline.OnNegativeDecision(isScheduledCardStep), pipeline.DevNull[card.BriefCard]())
	collectScheduledCardsStep := pipeline.WithFilter[*card.Card, *card.Card](p, limitScheduledCardsStep, scheduledCardsCollector)
	countScheduledCardsStep := pipeline.WithFilter[*card.Card, int](p, collectScheduledCardsStep, pipeline.NewCounter[*card.Card]())

	cardsCountCalculationStep := pipeline.WithAggregator[int](p, countNewCardsStep, countScheduledCardsStep, pipeline.SumCalculator())
	isCardLimitExceededStep := pipeline.WithCondition[int](p, cardsCountCalculationStep, pipeline.NewPredicateCondition(c.isCardLimitExceeded))
	pipeline.WithAcceptor(p, pipeline.OnPositiveDecision(isCardLimitExceededStep), pipeline.StopProcessing[int](p))
	pipeline.WithAcceptor(p, pipeline.OnNegativeDecision(isCardLimitExceededStep), pipeline.DevNull[int]())

	err := p.SyncRun()
	if err != nil {
		return nil, err
	}

	return &DailyCards{
		NewCards:       newCardsCollector.Items,
		ScheduledCards: scheduledCardsCollector.Items,
	}, nil
}

func (c *DailyCardsCollector) isCardLimitExceeded(count int) bool {
	return count >= c.config.NewCardsCount+c.config.ScheduledCardsCount
}
