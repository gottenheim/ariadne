package card_test

import (
	"testing"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/datetime"
)

func TestIsCardRemindedToday_IfNoRemindActivityCreated(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()
	activityChain := createTestActivityChain(learnCard)

	isRemindedToday, err := card.IsCardRemindedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isRemindedToday {
		t.Fatal("Card should not be reminded today, because there's no reminder activity")
	}
}

func TestIsCardRemindedToday_IfRemindActivityScheduledToTomorrow(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()
	activityChain := createTestActivityChain(learnCard, remindCard|remindCardScheduledToTomorrow)

	isRemindedToday, err := card.IsCardRemindedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isRemindedToday {
		t.Fatal("Card should not be reminded today, because reminder activity scheduled to tomorrow")
	}
}

func TestIsCardRemindedToday_IfRemindActivityScheduledToToday_ButNotExecuted(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()
	activityChain := createTestActivityChain(learnCard, remindCard|remindCardScheduledToToday)

	isRemindedToday, err := card.IsCardRemindedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isRemindedToday {
		t.Fatal("Card should not be reminded today, because reminder activity scheduled to today but not executed")
	}
}

func TestIsCardRemindedToday_IfRemindActivityScheduledToToday_AndExecuted(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()
	activityChain := createTestActivityChain(learnCard, remindCard|remindCardScheduledToToday|cardExecutedToday)

	isRemindedToday, err := card.IsCardRemindedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if !isRemindedToday {
		t.Fatal("Card should be reminded today, because reminder activity scheduled to today and executed")
	}
}

func TestIsCardRemindedToday_IfRemindActivityScheduledToYesterday_AndExecutedToday(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()
	activityChain := createTestActivityChain(learnCard, remindCard|remindCardScheduledToYesterday|cardExecutedToday)

	isRemindedToday, err := card.IsCardRemindedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if !isRemindedToday {
		t.Fatal("Card should be reminded today, because reminder activity scheduled to yesterday and executed today")
	}
}

func TestIsCardRemindedToday_IfRemindActivityScheduledToYesterday_AndExecutedYesterday(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()
	activityChain := createTestActivityChain(learnCard, remindCard|remindCardScheduledToYesterday|cardExecutedYesterday)

	isRemindedToday, err := card.IsCardRemindedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isRemindedToday {
		t.Fatal("Card should not be reminded today, because reminder activity scheduled to yesterday and executed yesterday")
	}
}

func TestIsCardRemindedToday_IfRemindActivityScheduledToToday_AndExecutedToday_AndNewRemindActivityHasAlreadyScheduled(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()
	activityChain := createTestActivityChain(learnCard, remindCard|remindCardScheduledToToday|cardExecutedToday, remindCard|remindCardScheduledToTomorrow)

	isRemindedToday, err := card.IsCardRemindedToday(timeSource, activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if !isRemindedToday {
		t.Fatal("Card should be reminded today, because reminder activity scheduled to today and executed today")
	}
}
