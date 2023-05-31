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

	if reminder.EasinessFactor() != 2.6 {
		t.Fatal("Reminder should increase one day interval by 0.1")
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

	if reminder.EasinessFactor() != 2.7 {
		t.Fatal("Reminder should increase one day interval by 0.1")
	}
}

func TestCardWorkflow_TransitionToFifthGradeThreeTimes_ShouldGenerateSixteenDayInterval(t *testing.T) {
	card := card.NewFakeCard().
		WithSection("languages/cpp").
		WithEntry("1").
		WithActivityChain(card.GenerateActivityChain(card.LearnCard)).
		Build()

	reminder := getReminderAfterChoosingGradeNTimes(t, card, 5, 3)

	// round(6*2.6)
	if reminder.Interval() != 16 {
		t.Fatal("Reminder should have six days interval")
	}

	if reminder.EasinessFactor() < 2.79 || reminder.EasinessFactor() > 2.81 {
		t.Fatal("Reminder should increase one day interval by 0.1")
	}
}

func TestCardWorkflow_TransitionToThirdGradeThreeTimes_ShouldGenerateThirteenDayInterval(t *testing.T) {
	card := card.NewFakeCard().
		WithSection("languages/cpp").
		WithEntry("1").
		WithActivityChain(card.GenerateActivityChain(card.LearnCard)).
		Build()

	reminder := getReminderAfterChoosingGradeNTimes(t, card, 3, 3)

	// round(6*2.08)
	if reminder.Interval() != 13 {
		t.Fatal("Reminder should have thirteen days interval")
	}

	// step 1. q=3, 2.5+(0.1-(5-q)*(0.08+(5-q)*0.02))=2.36
	// step 2. q=3, 2.36+(0.1-(5-q)*(0.08+(5-q)*0.02))=2.22
	// step 3. q=3, 2.22+(0.1-(5-q)*(0.08+(5-q)*0.02))=2.08
	if reminder.EasinessFactor() < 2.07 || reminder.EasinessFactor() > 2.09 {
		t.Fatal("Wrong easiness factor according to formula")
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
