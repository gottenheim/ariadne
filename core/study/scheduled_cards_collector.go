package study

import (
	"context"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/libraries/datetime"
	"github.com/gottenheim/ariadne/libraries/pipeline"
)

type scheduledCardsCollector struct {
	timeSource       datetime.TimeSource
	cardRepo         card.CardRepository
	config           *DailyCardsConfig
	scheduledCards   []*card.Card
	hotCardsToRevise []*card.Card
}

func CollectScheduledCards(timeSource datetime.TimeSource, cardRepo card.CardRepository, config *DailyCardsConfig) *scheduledCardsCollector {
	return &scheduledCardsCollector{
		timeSource: timeSource,
		cardRepo:   cardRepo,
		config:     config,
	}
}

func (f *scheduledCardsCollector) Run(ctx context.Context, input <-chan card.BriefCard, output chan<- card.BriefCard) error {
	var scheduledToRemindToday, remindedToday, hotCardsToRevise []*card.Card

	for {
		briefCard, ok := <-input
		if !ok {
			break
		}

		isScheduledToRemindOrRemindedToday, err := f.isCardScheduledToRemindOrRemindedToday(briefCard)
		if err != nil {
			return err
		}

		if !isScheduledToRemindOrRemindedToday {
			pipeline.WriteToChannel[card.BriefCard](ctx, output, briefCard)
			continue
		}

		crd, err := f.cardRepo.Get(briefCard.Section, briefCard.Entry)

		if err != nil {
			return err
		}

		isScheduledToRemindToday, err := card.IsCardScheduledToRemindToday(f.timeSource, briefCard.Activities)

		if err != nil {
			return err
		}

		if isScheduledToRemindToday {
			isRemindedToday, err := card.IsCardRemindedToday(f.timeSource, crd.Activities())

			if err != nil {
				return err
			}

			if !isRemindedToday {
				scheduledToRemindToday = append(scheduledToRemindToday, crd)
			} else {
				hotCardsToRevise = append(hotCardsToRevise, crd)
			}
		} else {
			remindedToday = append(remindedToday, crd)
		}
	}

	scheduledCardsRemaining := f.config.ScheduledCardsCount - len(remindedToday) - len(hotCardsToRevise)
	if scheduledCardsRemaining < 0 {
		scheduledCardsRemaining = 0
	}

	if len(scheduledToRemindToday) < scheduledCardsRemaining {
		scheduledCardsRemaining = len(scheduledToRemindToday)
	}

	f.scheduledCards = scheduledToRemindToday[0:scheduledCardsRemaining]
	f.hotCardsToRevise = hotCardsToRevise

	return nil
}

func (f *scheduledCardsCollector) isCardScheduledToRemindOrRemindedToday(briefCard card.BriefCard) (bool, error) {
	isScheduledToRemindToday, err := card.IsCardScheduledToRemindToday(f.timeSource, briefCard.Activities)

	if err != nil {
		return false, err
	}

	if isScheduledToRemindToday {
		return true, nil
	}

	isRemindedToday, err := card.IsCardRemindedToday(f.timeSource, briefCard.Activities)

	if err != nil {
		return false, err
	}

	return isRemindedToday, nil
}
