package study

import (
	"math/rand"
	"sort"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/libraries/datetime"
)

type DailyCardType int

const (
	HotDailyCard DailyCardType = iota
	NewDailyCard
	ScheduledDailyCard
)

type SelectedDailyCard struct {
	CardType DailyCardType
	Card     *card.Card
}

type StudyProgress struct {
	HotCardsleft       int
	NewCardsLeft       int
	ScheduledCardsLeft int
}

type DailyCardsIterator struct {
	timeSource       datetime.TimeSource
	newCards         []*card.Card
	scheduledCards   []*card.Card
	hotCardsToRevise []*card.Card
}

func NewDailyCardsIterator(timeSource datetime.TimeSource, dailyCards *DailyCards) *DailyCardsIterator {
	return &DailyCardsIterator{
		timeSource:       timeSource,
		newCards:         dailyCards.NewCards,
		scheduledCards:   dailyCards.ScheduledCards,
		hotCardsToRevise: sortCardsByTime(timeSource, dailyCards.HotCardsToRevise),
	}
}

func (i *DailyCardsIterator) Next() (*SelectedDailyCard, error) {
	randomCardsTotal := len(i.newCards) + len(i.scheduledCards)

	if len(i.hotCardsToRevise) > 0 {
		firstCard := i.hotCardsToRevise[0]
		canReturnHotCard := randomCardsTotal == 0

		if !canReturnHotCard {
			cardTime, err := card.GetTimeToRemindToday(i.timeSource, firstCard.Activities())
			if err != nil {
				return nil, err
			}

			canReturnHotCard = i.timeSource.Now().After(cardTime)
		}

		if canReturnHotCard {
			i.hotCardsToRevise = i.hotCardsToRevise[1:]
			return &SelectedDailyCard{
					CardType: HotDailyCard,
					Card:     firstCard},
				nil
		}
	}

	if randomCardsTotal == 0 {
		return nil, nil
	}

	randomIndex := rand.Int() % randomCardsTotal

	if randomIndex < len(i.newCards) {
		newCard := i.newCards[0]
		i.newCards = i.newCards[1:]
		return &SelectedDailyCard{
			CardType: NewDailyCard,
			Card:     newCard,
		}, nil
	} else {
		randomIndex -= len(i.newCards)
		scheduledCard := i.scheduledCards[randomIndex]
		i.scheduledCards = append(i.scheduledCards[0:randomIndex], i.scheduledCards[randomIndex+1:]...)
		return &SelectedDailyCard{
			CardType: ScheduledDailyCard,
			Card:     scheduledCard,
		}, nil
	}
}

func (i *DailyCardsIterator) AddHotCardToRevise(crd *card.Card) {
	i.hotCardsToRevise = append(i.hotCardsToRevise, crd)
}

func (i *DailyCardsIterator) GetStudyProgress() *StudyProgress {
	return &StudyProgress{
		HotCardsleft:       len(i.hotCardsToRevise),
		NewCardsLeft:       len(i.newCards),
		ScheduledCardsLeft: len(i.scheduledCards),
	}
}

func sortCardsByTime(timeSource datetime.TimeSource, hotCards []*card.Card) []*card.Card {
	sort.Slice(hotCards, func(i, j int) bool {
		leftTime, _ := card.GetTimeToRemindToday(timeSource, hotCards[i].Activities())
		rightTime, _ := card.GetTimeToRemindToday(timeSource, hotCards[j].Activities())
		return leftTime.Before(rightTime)
	})
	return hotCards
}
