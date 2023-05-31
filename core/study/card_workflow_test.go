package study_test

import (
	"testing"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/core/study"
	"github.com/gottenheim/ariadne/libraries/datetime"
)

func TestCardWorkflow_SingleTransitionToFifthGrade_ShouldGenerateOneDayInterval(t *testing.T) {
	card := card.NewFakeCard().
		WithSection("languages/cpp").
		WithEntry("1").
		WithActivityChain(card.GenerateActivityChain(card.LearnCard)).
		Build()

	reminder := getReminderAfterChoosingGradeNTimes(t, card, 5, 1)

	if reminder.Interval() != 1 {
		t.Fatal("Reminder should have one day interval")
	}
}

func TestCardWorkflow_TransitionToFifthGradeTwice_ShouldGenerateSixDayInterval(t *testing.T) {
	card := card.NewFakeCard().
		WithSection("languages/cpp").
		WithEntry("1").
		WithActivityChain(card.GenerateActivityChain(card.LearnCard)).
		Build()

	reminder := getReminderAfterChoosingGradeNTimes(t, card, 5, 2)

	if reminder.Interval() != 6 {
		t.Fatal("Reminder should have six days interval")
	}
}

func getReminderAfterChoosingGradeNTimes(t *testing.T, c *card.Card, grade int, times int) *card.RemindCardActivity {
	timeSource := datetime.NewFakeTimeSource()

	for i := 0; i < times; i++ {
		workflow := study.NewCardWorkflow(timeSource, c)

		nextStates, err := workflow.GetNextStates()

		if err != nil {
			t.Fatal(err)
		}

		err = workflow.TransitTo(nextStates[grade])

		if err != nil {
			t.Fatal(err)
		}
	}

	lastActivity := c.Activities()

	remindCard := lastActivity.(*card.RemindCardActivity)

	return remindCard
}
