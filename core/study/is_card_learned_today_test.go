package study_test

import (
	"testing"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/core/study"
	"github.com/gottenheim/ariadne/libraries/datetime"
)

func TestIsCardLearnedToday_IfLearnActivityIsNotExecuted(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()
	activityChain := card.GenerateActivityChain(card.LearnCard)

	isLearnedToday, err := study.IsCardLearnedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isLearnedToday {
		t.Fatal("Card should not be learned today, because learn activity has not been executed yet")
	}
}

func TestIsCardLearnedToday_IfLearnActivityIsNotExecuted_AndRemindActivityInTheEndOfChain(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()
	activityChain := card.GenerateActivityChain(card.LearnCard, card.RemindCard)

	isLearnedToday, err := study.IsCardLearnedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isLearnedToday {
		t.Fatal("Card should not be learned today, because learn activity has not been executed yet")
	}
}

func TestIsCardLearnedToday_IfLearnActivityHasBeenExecutedYesterday(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()

	activityChain := card.GenerateActivityChain(card.LearnCard | card.CardExecutedYesterday)

	isLearnedToday, err := study.IsCardLearnedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isLearnedToday {
		t.Fatal("Card should not be learned today, because learn activity was executed tomorrow")
	}
}

func TestIsCardLearnedToday_IfLearnActivityHasBeenExecutedToday(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()

	activityChain := card.GenerateActivityChain(card.LearnCard | card.CardExecutedToday)

	isLearnedToday, err := study.IsCardLearnedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if !isLearnedToday {
		t.Fatal("Card should be learned today, because learn activity was executed today")
	}
}

func TestIsCardLearnedToday_IfLearnActivityHasBeenExecutedToday_AndSomeRemindActivitiesAddedAfterIt(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()

	activityChain := card.GenerateActivityChain(card.LearnCard|card.CardExecutedToday, card.RemindCard|card.RemindCardScheduledToToday, card.RemindCard|card.RemindCardScheduledToTomorrow)

	isLearnedToday, err := study.IsCardLearnedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if !isLearnedToday {
		t.Fatal("Card should be learned today, because learn activity was executed today")
	}
}
