package study_test

import (
	"testing"
	"time"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/core/study"
	"github.com/gottenheim/ariadne/infra/interactor"
	"github.com/gottenheim/ariadne/libraries/datetime"
)

type isCardScheduledToRemindAt struct {
	timeSource datetime.TimeSource
	day        time.Time
	result     bool
}

func (s *isCardScheduledToRemindAt) OnLearnCard(learn *card.LearnCardActivity) error {
	return nil
}

func (s *isCardScheduledToRemindAt) OnRemindCard(remind *card.RemindCardActivity) error {
	s.result = !remind.IsExecuted() && datetime.IsSameDay(s.timeSource, remind.ScheduledTo(), s.day)

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

func studyCards(t *testing.T, timeSource datetime.TimeSource, config *study.DailyCardsConfig, cards []*card.Card, chooseState study.ChooseStateFunc) {
	cardRepo := card.NewFakeCardRepository(cards...)

	cardEmitter := &fakeCardEmitter{
		briefCards: card.ExtractBriefCards(cards),
	}

	session := study.NewSession(timeSource, cardRepo, interactor.NewFakeUserInteractor())

	err := session.Run(config, cardEmitter, chooseState)

	if err != nil {
		t.Fatal(err)
	}
}

func rememberAllCardsWell(crd *card.Card, states []*study.CardState) (*study.CardState, error) {
	return getStateByGrade(states, Good), nil
}

func forgotNCards(cardsToForget int) study.ChooseStateFunc {
	cardsToForgetLeft := cardsToForget
	return func(crd *card.Card, states []*study.CardState) (*study.CardState, error) {
		if cardsToForgetLeft > 0 {
			cardsToForgetLeft--
			return getStateByGrade(states, Again), nil
		} else {
			return getStateByGrade(states, Good), nil
		}
	}
}

func TestSession_ShouldStudyTenNewCards_AndScheduleThemToTomorrow(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 100, card.LearnCard)).
		Generate()

	config := &study.DailyCardsConfig{
		NewCardsCount:       10,
		ScheduledCardsCount: 0,
	}

	timeSource := datetime.NewFakeTimeSource()

	studyCards(t, timeSource, config, cards, rememberAllCardsWell)

	tomorrow := datetime.GetToday(timeSource).AddDate(0, 0, 1)

	assertMatchingCardsCount(t, cards, func(c *card.Card) (bool, error) {
		return IsCardScheduledToRemindAt(c, tomorrow)
	}, 10)
}

func TestSession_ShouldNotStudyMoreNewCards_IfAllNewCardsForTodayAreAlreadyLearned(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 100, card.LearnCard)).
		Generate()

	config := &study.DailyCardsConfig{
		NewCardsCount:       10,
		ScheduledCardsCount: 0,
	}

	timeSource := datetime.NewFakeTimeSource()

	studyCards(t, timeSource, config, cards, rememberAllCardsWell)
	studyCards(t, timeSource, config, cards, rememberAllCardsWell)

	tomorrow := datetime.GetToday(timeSource).AddDate(0, 0, 1)

	assertMatchingCardsCount(t, cards, func(c *card.Card) (bool, error) {
		return IsCardScheduledToRemindAt(c, tomorrow)
	}, 10)

	assertMatchingCardsCount(t, cards, func(c *card.Card) (bool, error) {
		return card.IsCardNew(c.Activities())
	}, 90)
}

func TestSession_ShouldRescheduleCardsToSixDaysAhead_IfCardsCanBeRemembered(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 100, card.LearnCard)).
		Generate()

	timeSource := datetime.NewFakeTimeSource()

	studyCards(t, timeSource, &study.DailyCardsConfig{
		NewCardsCount:       10,
		ScheduledCardsCount: 0,
	}, cards, rememberAllCardsWell)

	timeSource.MoveNow(1 * study.Day)

	studyCards(t, timeSource, &study.DailyCardsConfig{
		NewCardsCount:       0,
		ScheduledCardsCount: 10,
	}, cards, rememberAllCardsWell)

	sixDaysAhead := datetime.GetToday(timeSource).AddDate(0, 0, 6)

	assertMatchingCardsCount(t, cards, func(c *card.Card) (bool, error) {
		return IsCardScheduledToRemindAt(c, sixDaysAhead)
	}, 10)
}

func TestSession_ShouldRescheduleForgottenCardsToTomorrow(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 100, card.LearnCard)).
		Generate()

	timeSource := datetime.NewFakeTimeSource()

	studyCards(t, timeSource, &study.DailyCardsConfig{
		NewCardsCount:       10,
		ScheduledCardsCount: 0,
	}, cards, rememberAllCardsWell)

	timeSource.MoveNow(1 * study.Day)

	studyCards(t, timeSource, &study.DailyCardsConfig{
		NewCardsCount:       0,
		ScheduledCardsCount: 10,
	}, cards, forgotNCards(3))

	tomorrow := datetime.GetToday(timeSource).AddDate(0, 0, 1)

	assertMatchingCardsCount(t, cards, func(c *card.Card) (bool, error) {
		return IsCardScheduledToRemindAt(c, tomorrow)
	}, 3)

	sixDaysAhead := datetime.GetToday(timeSource).AddDate(0, 0, 6)

	assertMatchingCardsCount(t, cards, func(c *card.Card) (bool, error) {
		return IsCardScheduledToRemindAt(c, sixDaysAhead)
	}, 7)
}

func TestSession_ShouldNotStudyMoreScheduledCards_IfAllScheduledCardsForTodayAreAlreadyReminded(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 100, card.LearnCard)).
		Generate()

	timeSource := datetime.NewFakeTimeSource()

	studyCards(t, timeSource, &study.DailyCardsConfig{
		NewCardsCount:       10,
		ScheduledCardsCount: 0,
	}, cards, rememberAllCardsWell)

	timeSource.MoveNow(1 * study.Day)

	studyCards(t, timeSource, &study.DailyCardsConfig{
		NewCardsCount:       0,
		ScheduledCardsCount: 5,
	}, cards, rememberAllCardsWell)

	studyCards(t, timeSource, &study.DailyCardsConfig{
		NewCardsCount:       0,
		ScheduledCardsCount: 5,
	}, cards, rememberAllCardsWell)

	sixDaysAhead := datetime.GetToday(timeSource).AddDate(0, 0, 6)

	assertMatchingCardsCount(t, cards, func(c *card.Card) (bool, error) {
		return IsCardScheduledToRemindAt(c, sixDaysAhead)
	}, 5)
}

func TestSession_ShouldStudyNewCardsAndRepeatScheduledInTheSameSession(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 100, card.LearnCard)).
		Generate()

	config := &study.DailyCardsConfig{
		NewCardsCount:       10,
		ScheduledCardsCount: 10,
	}

	timeSource := datetime.NewFakeTimeSource()

	studyCards(t, timeSource, config, cards, rememberAllCardsWell)

	timeSource.MoveNow(1 * study.Day)

	studyCards(t, timeSource, config, cards, rememberAllCardsWell)

	tomorrow := datetime.GetToday(timeSource).AddDate(0, 0, 1)

	assertMatchingCardsCount(t, cards, func(c *card.Card) (bool, error) {
		return IsCardScheduledToRemindAt(c, tomorrow)
	}, 10)

	sixDaysAhead := datetime.GetToday(timeSource).AddDate(0, 0, 6)

	assertMatchingCardsCount(t, cards, func(c *card.Card) (bool, error) {
		return IsCardScheduledToRemindAt(c, sixDaysAhead)
	}, 10)
}
