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
	NewCards         []*card.Card
	ScheduledCards   []*card.Card
	HotCardsToRevise []*card.Card
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
	newCardsCollector := CollectNewCards(c.timeSource, c.cardRepo, c.config)
	scheduledCardsCollector := CollectScheduledCards(c.timeSource, c.cardRepo, c.config)

	p := pipeline.New()
	cardEmissionStep := pipeline.NewEmitter(p, c.cardEmitter)
	collectNewCardsStep := pipeline.WithFilter[card.BriefCard, card.BriefCard](p, cardEmissionStep, newCardsCollector)
	collectScheduledCardsStep := pipeline.WithFilter[card.BriefCard, card.BriefCard](p, collectNewCardsStep, scheduledCardsCollector)
	pipeline.WithAcceptor[card.BriefCard](p, collectScheduledCardsStep, pipeline.Skip[card.BriefCard]())

	err := p.SyncRun()
	if err != nil {
		return nil, err
	}

	return &DailyCards{
		NewCards:         newCardsCollector.newCards,
		ScheduledCards:   scheduledCardsCollector.scheduledCards,
		HotCardsToRevise: append(newCardsCollector.hotCardsToRevise, scheduledCardsCollector.hotCardsToRevise...),
	}, nil
}
