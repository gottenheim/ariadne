package study_test

import (
	"testing"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/core/study"
	"github.com/gottenheim/ariadne/libraries/datetime"
)

func TestIsCardRemindedToday_IfNoRemindActivityCreated(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()
	activityChain := card.GenerateActivityChain(card.LearnCard)

	isRemindedToday, err := study.IsCardRemindedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isRemindedToday {
		t.Fatal("Card should not be reminded today, because there's no reminder activity")
	}
}

func TestIsCardRemindedToday_IfRemindActivityScheduledToTomorrow(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()
	activityChain := card.GenerateActivityChain(card.LearnCard, card.RemindCard|card.RemindCardScheduledToTomorrow)

	isRemindedToday, err := study.IsCardRemindedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isRemindedToday {
		t.Fatal("Card should not be reminded today, because reminder activity scheduled to tomorrow")
	}
}

func TestIsCardRemindedToday_IfRemindActivityScheduledToToday_ButNotExecuted(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()
	activityChain := card.GenerateActivityChain(card.LearnCard, card.RemindCard|card.RemindCardScheduledToToday)

	isRemindedToday, err := study.IsCardRemindedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isRemindedToday {
		t.Fatal("Card should not be reminded today, because reminder activity scheduled to today but not executed")
	}
}

func TestIsCardRemindedToday_IfRemindActivityScheduledToToday_AndExecuted(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()
	activityChain := card.GenerateActivityChain(card.LearnCard, card.RemindCard|card.RemindCardScheduledToToday|card.CardExecutedToday)

	isRemindedToday, err := study.IsCardRemindedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if !isRemindedToday {
		t.Fatal("Card should be reminded today, because reminder activity scheduled to today and executed")
	}
}

func TestIsCardRemindedToday_IfRemindActivityScheduledToYesterday_AndExecutedToday(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()
	activityChain := card.GenerateActivityChain(card.LearnCard, card.RemindCard|card.RemindCardScheduledToYesterday|card.CardExecutedToday)

	isRemindedToday, err := study.IsCardRemindedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if !isRemindedToday {
		t.Fatal("Card should be reminded today, because reminder activity scheduled to yesterday and executed today")
	}
}

func TestIsCardRemindedToday_IfRemindActivityScheduledToYesterday_AndExecutedYesterday(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()
	activityChain := card.GenerateActivityChain(card.LearnCard, card.RemindCard|card.RemindCardScheduledToYesterday|card.CardExecutedYesterday)

	isRemindedToday, err := study.IsCardRemindedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isRemindedToday {
		t.Fatal("Card should not be reminded today, because reminder activity scheduled to yesterday and executed yesterday")
	}
}

func TestIsCardRemindedToday_IfRemindActivityScheduledToToday_AndExecutedToday_AndNewRemindActivityHasAlreadyScheduled(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()
	activityChain := card.GenerateActivityChain(card.LearnCard, card.RemindCard|card.RemindCardScheduledToToday|card.CardExecutedToday, card.RemindCard|card.RemindCardScheduledToTomorrow)

	isRemindedToday, err := study.IsCardRemindedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if !isRemindedToday {
		t.Fatal("Card should be reminded today, because reminder activity scheduled to today and executed today")
	}
}
