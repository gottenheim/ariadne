package card_test

import (
	"testing"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/pipeline"
)

type cardGenerator struct {
	events pipeline.FilterEvents
	cards  []*card.KeyWithActivities
}

func newCardGenerator(events pipeline.FilterEvents, cards ...*card.KeyWithActivities) *cardGenerator {
	events.OnStart()
	return &cardGenerator{
		events: events,
		cards:  cards,
	}
}

func (g *cardGenerator) Run(output chan<- *card.KeyWithActivities) {
	defer func() {
		close(output)
		g.events.OnFinish()
	}()

	for _, card := range g.cards {
		output <- card
	}
}

type cardFilteringResult struct {
	cards []*card.Card
}

type cardAccumulator struct {
	events pipeline.FilterEvents
	result *cardFilteringResult
}

func newCardAccumulator(events pipeline.FilterEvents, result *cardFilteringResult) *cardAccumulator {
	events.OnStart()
	return &cardAccumulator{
		events: events,
		result: result,
	}
}

func (g *cardAccumulator) Run(input <-chan *card.Card, output chan<- interface{}) {
	defer func() {
		g.events.OnFinish()
	}()

	for {
		card, ok := <-input
		if !ok {
			break
		}
		g.result.cards = append(g.result.cards, card)
	}
}

func TestNewCardFilter(t *testing.T) {
	events := &pipeline.WaitGroupEventHandler{}

	cards := card.NewBatchCardGenerator().
		WithNewCards(80).
		WithLearnedCards(90).
		WithRemindedCards(100).
		Generate()

	keysWithActivities := card.ExtractKeysWithActivities(cards)

	cardRepo := card.NewFakeCardRepository(cards...)

	filteringResult := &cardFilteringResult{}

	cardPipeline := pipeline.Join[*card.Card, interface{}](
		pipeline.Join[*card.KeyWithActivities](
			pipeline.New[*card.KeyWithActivities](
				newCardGenerator(events, keysWithActivities...)),
			card.NewCardFilter(events, cardRepo)),
		newCardAccumulator(events, filteringResult))

	cardPipeline.Run()

	events.Wait()

	if len(filteringResult.cards) != 80 {
		t.Fatal("Filter should find ten cards")
	}
}
