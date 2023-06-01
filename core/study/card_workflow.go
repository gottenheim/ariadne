package study

import (
	"time"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/libraries/datetime"
)

type CardState struct {
	Name             string
	Description      string
	Grade            int
	EasinessFactor   float64
	RepetitionNumber int
	Interval         time.Duration
}

const (
	Day time.Duration = 24 * time.Hour
)

type CardWorkflow struct {
	timeSource datetime.TimeSource
	card       *card.Card
}

func NewCardWorkflow(timeSource datetime.TimeSource, card *card.Card) *CardWorkflow {
	return &CardWorkflow{
		timeSource: timeSource,
		card:       card,
	}
}

func (w *CardWorkflow) GetNextStates() ([]*CardState, error) {
	repetitionParams, err := GetCardRepetitionParams(w.card.Activities())

	if err != nil {
		return nil, err
	}

	nextStates := []*CardState{
		{
			Name:        "Again",
			Description: "Complete failure to recall the information",
			Grade:       0,
		},
		{
			Name:        "Hard",
			Description: "Correct response, but required significant effort to recall",
			Grade:       3,
		},
		{
			Name:        "Good",
			Description: "Correct response, after some hesitation",
			Grade:       4,
		},
		{
			Name:        "Easy",
			Description: "Correct response with perfect recall",
			Grade:       5,
		},
	}

	for _, nextState := range nextStates {
		nextState.RepetitionNumber = repetitionParams.RepetitionNumber + 1

		if nextState.Grade >= 3 {
			if repetitionParams.RepetitionNumber == 0 {
				nextState.Interval = 10 * time.Minute
			} else if repetitionParams.RepetitionNumber == 1 {
				nextState.Interval = 1 * Day
			} else if repetitionParams.RepetitionNumber == 2 {
				nextState.Interval = 6 * Day
			} else {
				intervalInDays := repetitionParams.Interval.Hours() / 24
				nextState.Interval = time.Duration(int(intervalInDays*repetitionParams.EasinessFactor)) * Day
			}
		} else {
			nextState.RepetitionNumber = 1
			nextState.Interval = 10 * time.Minute
		}

		nextState.EasinessFactor = repetitionParams.EasinessFactor + (0.1 - float64(5-nextState.Grade)*(0.08+float64(5-nextState.Grade)*0.02))
		if nextState.EasinessFactor < 1.3 {
			nextState.EasinessFactor = 1.3
		}
	}

	return nextStates, nil
}

func (w *CardWorkflow) TransitTo(state *CardState) error {
	now := w.timeSource.Today()
	err := MarkLastActivityAsExecuted(w.card.Activities(), now)
	if err != nil {
		return err
	}

	nextReminder := card.CreateRemindCardActivity(w.card.Activities())
	nextReminder.SetEasinessFactor(state.EasinessFactor)
	nextReminder.SetRepetitionNumber(state.RepetitionNumber)
	nextReminder.SetInterval(state.Interval)
	nextReminder.ScheduleTo(now.Add(state.Interval))

	w.card.SetActivities(nextReminder)

	return nil
}
