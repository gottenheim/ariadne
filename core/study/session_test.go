package study_test

import (
	"testing"
	"time"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/core/study"
	"github.com/gottenheim/ariadne/libraries/datetime"
)

type isCardScheduledToRemindAt struct {
	day    time.Time
	result bool
}

func (s *isCardScheduledToRemindAt) OnLearnCard(learn *card.LearnCardActivity) error {
	return nil
}

func (s *isCardScheduledToRemindAt) OnRemindCard(remind *card.RemindCardActivity) error {
	s.result = !remind.IsExecuted() && remind.ScheduledTo() == s.day

	if s.result {
		return nil
	}

	if remind.PreviousActivity() == nil {
		return nil
	}

	return remind.PreviousActivity().Accept(s)
}

func IsCardScheduledToRemindAt(crd *card.Card, day time.Time) (bool, error) {
	isScheduledToRemind := &isCardScheduledToRemindAt{
		day:    day,
		result: false,
	}
	err := crd.Activities().Accept(isScheduledToRemind)
	if err != nil {
		return false, err
	}
	return isScheduledToRemind.result, nil
}

func assertMatchingCardsCount(t *testing.T, cards []*card.Card, matchingFunc func(*card.Card) (bool, error), expectedCount int) {
	countMatching := 0
	for _, crd := range cards {
		isMatching, err := matchingFunc(crd)

		if err != nil {
			t.Fatal(err)
		}

		if isMatching {
			countMatching++
		}
	}
	if countMatching != expectedCount {
		t.Errorf("Expected %d matching cards but found %d", expectedCount, countMatching)
	}
}

func TestSession_ShouldStudyTenCards_AndScheduleThemToTomorrow(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()

	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Cards scheduled to today", 100, card.LearnCard)).
		Generate()

	cardRepo := card.NewFakeCardRepository(cards...)

	config := &study.DailyCardsConfig{
		NewCardsCount:       10,
		ScheduledCardsCount: 0,
	}

	cardEmitter := &fakeCardEmitter{
		briefCards: card.ExtractBriefCards(cards),
	}

	session := study.NewSession(timeSource, cardRepo)

	err := session.Run(config, cardEmitter, func(states []*study.CardState) (*study.CardState, error) {
		return getStateByGrade(states, Good), nil
	})

	if err != nil {
		t.Fatal(err)
	}

	tomorrow := timeSource.Today().AddDate(0, 0, 1)

	assertMatchingCardsCount(t, cards, func(c *card.Card) (bool, error) {
		return IsCardScheduledToRemindAt(c, tomorrow)
	}, 10)
}
