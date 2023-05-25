package study_test

import (
	"testing"

	"github.com/gottenheim/ariadne/card"
)

func TestStudyCardFilter(t *testing.T) {
	// cards := generateCardsToStudy(&cardGenerationConfig{
	// 	newCards: 10,
	// })
	// repo := card.NewFakeCardRepository(cards...)
	// timeSource := datetime.NewFakeTimeSource()

	// filter := study.NewCardFilter(repo, timeSource, 5, 0)

}

type cardGenerationConfig struct {
	newCards int
}

func generateCardsToStudy(cfg *cardGenerationConfig) []*card.Card {
	var cards []*card.Card

	orderNum := 1

	for i := 0; i < cfg.newCards; i++ {
		c := card.NewCard([]string{"books", "cpp"}, orderNum, []card.CardArtifact{})
		c.SetActivities(card.GenerateActivityChain(card.LearnCard))
		orderNum++
	}

	return cards
}
