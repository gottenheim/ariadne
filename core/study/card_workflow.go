package study

import (
	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/libraries/datetime"
)

type CardState struct {
	Name             string
	Grade            int
	EasinessFactor   float32
	RepetitionNumber int
	Interval         int
}

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
			Name:  "Complete failure to recall the information",
			Grade: 0,
		},
		{
			Name:  "Incorrect response, but upon seeing the correct answer it felt familiar",
			Grade: 1,
		},
		{
			Name:  "Incorrect response, but upon seeing the correct answer it seemed easy to remember",
			Grade: 2,
		},
		{
			Name:  "Correct response, but required significant effort to recall",
			Grade: 3,
		},
		{
			Name:  "Correct response, after some hesitation",
			Grade: 4,
		},
		{
			Name:  "Correct response with perfect recall",
			Grade: 5,
		},
	}

	for _, nextState := range nextStates {
		nextState.RepetitionNumber = repetitionParams.repetitionNumber + 1

		if nextState.Grade >= 3 {
			if repetitionParams.repetitionNumber == 0 {
				nextState.Interval = 1
			} else if repetitionParams.repetitionNumber == 1 {
				nextState.Interval = 6
			} else {
				nextState.Interval = int(float32(repetitionParams.interval) * repetitionParams.easinessFactor)
			}
		} else {
			nextState.RepetitionNumber = 0
			nextState.Interval = 1
		}

		nextState.EasinessFactor = repetitionParams.easinessFactor + (0.1 - float32(5-nextState.Grade)*(0.08+float32(5-nextState.Grade)*0.02))
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
	nextReminder.ScheduleTo(now.AddDate(0, 0, state.Interval))

	w.card.SetActivities(nextReminder)

	return nil
}
