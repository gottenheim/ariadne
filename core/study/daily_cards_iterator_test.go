package study_test

import (
	"testing"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/core/study"
	"github.com/gottenheim/ariadne/libraries/datetime"
)

func TestDailyCardsIterator_ShouldReturnNewCards(t *testing.T) {

	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 1000, card.LearnCard)).
		Generate()

	timeSource := datetime.NewFakeTimeSource()

	it := study.NewDailyCardsIterator(timeSource, &study.DailyCards{
		NewCards: cards,
	})

	cardsCount := getCardsCount(t, it)

	if cardsCount != 1000 {
		t.Errorf("Iterator must return 1000 cards, but returned %d", cardsCount)
	}
}

func TestDailyCardsIterator_ShouldReturnScheduledCards(t *testing.T) {
	cards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Scheduled cards", 1000, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToToday)).
		Generate()

	timeSource := datetime.NewFakeTimeSource()

	it := study.NewDailyCardsIterator(timeSource, &study.DailyCards{
		ScheduledCards: cards,
	})

	cardsCount := getCardsCount(t, it)

	if cardsCount != 1000 {
		t.Errorf("Iterator must return 1000 cards, but returned %d", cardsCount)
	}
}

func TestDailyCardsIterator_ShouldReturnNewAndScheduledCards(t *testing.T) {
	newCards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 1000, card.LearnCard)).
		Generate()

	scheduledCards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Scheduled cards", 1000, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToToday)).
		Generate()

	timeSource := datetime.NewFakeTimeSource()

	it := study.NewDailyCardsIterator(timeSource, &study.DailyCards{
		NewCards:       newCards,
		ScheduledCards: scheduledCards,
	})

	cardsCount := getCardsCount(t, it)

	if cardsCount != 2000 {
		t.Errorf("Iterator must return 2000 cards, but returned %d", cardsCount)
	}
}

func TestDailyCardsIterator_ShouldReturnHotCardsFirst_IfTheirScheduleTimeIsExpired(t *testing.T) {
	newCards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 100, card.LearnCard)).
		Generate()

	scheduledCards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Scheduled cards", 100, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToToday)).
		Generate()

	hotCardsToRevise := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Hot cards to revise", 1, card.LearnCard|card.CardExecutedToday, card.RemindCard|card.RemindCardScheduledToFiveMinutesAgo)).
		Generate()

	timeSource := datetime.NewFakeTimeSource()

	it := study.NewDailyCardsIterator(timeSource, &study.DailyCards{
		NewCards:         newCards,
		ScheduledCards:   scheduledCards,
		HotCardsToRevise: hotCardsToRevise,
	})

	val, err := it.Next()

	if err != nil {
		t.Fatal(err)
	}

	if val.CardType != study.HotDailyCard {
		t.Error("Hot card to revise with expired time must be returned first")
	}
}

func TestDailyCardsIterator_ShouldNotReturnHotCardsFirst_IfTheirScheduleTimeIsNotExpiredYet(t *testing.T) {
	newCards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 100, card.LearnCard)).
		Generate()

	scheduledCards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Scheduled cards", 100, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToToday)).
		Generate()

	hotCardsToRevise := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Hot cards to revise", 1, card.LearnCard|card.CardExecutedToday, card.RemindCard|card.RemindCardScheduledToFiveMinutesAhead)).
		Generate()

	timeSource := datetime.NewFakeTimeSource()

	it := study.NewDailyCardsIterator(timeSource, &study.DailyCards{
		NewCards:         newCards,
		ScheduledCards:   scheduledCards,
		HotCardsToRevise: hotCardsToRevise,
	})

	val, err := it.Next()

	if err != nil {
		t.Fatal(err)
	}

	if val.CardType == study.HotDailyCard {
		t.Error("Hot card to revise with not expired time must not be returned first")
	}
}

func TestDailyCardsIterator_ShouldReturnHotCardsAfterAllOthers_EvenIfTheirTimeIsNotYetExpired(t *testing.T) {
	newCards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 100, card.LearnCard)).
		Generate()

	scheduledCards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Scheduled cards", 100, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToToday)).
		Generate()

	hotCardsToRevise := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Hot cards to revise", 1, card.LearnCard|card.CardExecutedToday, card.RemindCard|card.RemindCardScheduledToFiveMinutesAhead)).
		Generate()

	timeSource := datetime.NewFakeTimeSource()

	it := study.NewDailyCardsIterator(timeSource, &study.DailyCards{
		NewCards:         newCards,
		ScheduledCards:   scheduledCards,
		HotCardsToRevise: hotCardsToRevise,
	})

	cardsCount := getCardsCount(t, it)

	if cardsCount != 201 {
		t.Errorf("Iterator must return 201 cards, but returned %d", cardsCount)
	}
}

func TestDailyCardsIterator_ShouldReturnExpiredHotCardFirst_AndNotExpiredLast(t *testing.T) {
	newCards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 100, card.LearnCard)).
		Generate()

	scheduledCards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Scheduled cards", 100, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToToday)).
		Generate()

	hotCardsToRevise := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Hot cards to revise", 1, card.LearnCard|card.CardExecutedToday, card.RemindCard|card.RemindCardScheduledToFiveMinutesAgo)).
		WithCards(card.NewCardGenerationSpec("Hot cards to revise", 1, card.LearnCard|card.CardExecutedToday, card.RemindCard|card.RemindCardScheduledToFiveMinutesAhead)).
		Generate()

	timeSource := datetime.NewFakeTimeSource()

	it := study.NewDailyCardsIterator(timeSource, &study.DailyCards{
		NewCards:         newCards,
		ScheduledCards:   scheduledCards,
		HotCardsToRevise: hotCardsToRevise,
	})

	cards := getCards(t, it)

	if cards[0].Section() != "Hot cards to revise" {
		t.Error("Iterator must return expired hot card first")
	}

	if cards[len(cards)-1].Section() != "Hot cards to revise" {
		t.Error("Iterator must return not expired hot card last")
	}
}

func TestDailyCardsIterator_ShouldReturnNewAndScheduledCards_UsingUniformDistribution(t *testing.T) {
	newCards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("New cards", 1000, card.LearnCard)).
		Generate()

	scheduledCards := card.NewBatchCardGenerator().
		WithCards(card.NewCardGenerationSpec("Scheduled cards", 100, card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToToday)).
		Generate()

	timeSource := datetime.NewFakeTimeSource()

	it := study.NewDailyCardsIterator(timeSource, &study.DailyCards{
		NewCards:       newCards,
		ScheduledCards: scheduledCards,
	})

	cards := getCards(t, it)

	var runs []int
	runLen := 0

	for _, card := range cards {
		if card.Section() == "New cards" {
			runLen++
		} else {
			runs = append(runs, runLen)
			runLen = 0
		}
	}

	runSum := 0

	for _, runSize := range runs {
		runSum += runSize
	}

	avgRunLength := runSum / len(runs)

	if avgRunLength < 5 || avgRunLength > 15 {
		t.Error("Ideally iterator should return one scheduled card per 10 new cards")
	}
}

func getCardsCount(t *testing.T, it *study.DailyCardsIterator) int {
	cards := getCards(t, it)

	return len(cards)
}

func getCards(t *testing.T, it *study.DailyCardsIterator) []*card.Card {
	cardsCount := 0

	keyCardsMap := map[string]*card.Card{}

	var cards []*card.Card

	for ; ; cardsCount++ {
		val, err := it.Next()
		if err != nil {
			t.Fatal(err)
		}
		if val == nil {
			break
		}

		crd := val.Card

		key := crd.Section() + crd.Entry()

		_, ok := keyCardsMap[key]
		if ok {
			t.Fatal("Same card returned twice")
		}

		keyCardsMap[key] = crd

		cards = append(cards, crd)
	}

	return cards
}
