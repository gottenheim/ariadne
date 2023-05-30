package study_test

import (
	"testing"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/core/study"
	"github.com/gottenheim/ariadne/libraries/datetime"
)

func TestIsCardScheduledToRemindToday_IfNoReminderActivitiesExist(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()

	activityChain := card.GenerateActivityChain(card.LearnCard)

	isScheduledToRemind, err := study.IsCardScheduledToRemindToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isScheduledToRemind {
		t.Fatal("Card should not be scheduled to remind, because there's no reminder activity at all")
	}
}

func TestIsCardScheduledToRemindToday_IfReminderActivityExistsButScheduledToTomorrow(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()

	activityChain := card.GenerateActivityChain(card.LearnCard, card.RemindCard|card.RemindCardScheduledToTomorrow)

	isScheduledToRemind, err := study.IsCardScheduledToRemindToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isScheduledToRemind {
		t.Fatal("Card should not be scheduled to remind, because it's scheduled to tomorrow")
	}
}

func TestIsCardScheduledToRemindToday_IfReminderActivityExistsAndScheduledToToday(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()

	activityChain := card.GenerateActivityChain(card.LearnCard, card.RemindCard|card.RemindCardScheduledToToday)

	isScheduledToRemind, err := study.IsCardScheduledToRemindToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if !isScheduledToRemind {
		t.Fatal("Card should be scheduled to remind, because it's scheduled to today")
	}
}

func TestIsCardScheduledToRemindToday_IfReminderActivityExistsAndScheduledToToday_ButHasBeenAlreadyExecuted(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()

	activityChain := card.GenerateActivityChain(card.LearnCard, card.RemindCard|card.RemindCardScheduledToToday|card.CardExecutedToday)

	isScheduledToRemind, err := study.IsCardScheduledToRemindToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isScheduledToRemind {
		t.Fatal("Card should not be scheduled to remind, because it has been already executed")
	}
}

func TestIsCardScheduledToRemindToday_IfReminderActivityExistsAndScheduledToYesterday_ButHasNotBeenExecutedYet(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()

	activityChain := card.GenerateActivityChain(card.LearnCard, card.RemindCard|card.RemindCardScheduledToYesterday)

	isScheduledToRemind, err := study.IsCardScheduledToRemindToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if !isScheduledToRemind {
		t.Fatal("Card should be scheduled to remind, because it hasn't been executed yet")
	}
}

func TestIsCardScheduledToRemindToday_IfReminderActivityExistsAndScheduledToYesterday_AndHasBeenAlreadyExecutedToday(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()

	activityChain := card.GenerateActivityChain(card.LearnCard, card.RemindCard|card.RemindCardScheduledToYesterday|card.CardExecutedToday)

	isScheduledToRemind, err := study.IsCardScheduledToRemindToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isScheduledToRemind {
		t.Fatal("Card should not be scheduled to remind, because it has been already executed")
	}
}
