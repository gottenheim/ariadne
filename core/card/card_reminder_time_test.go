package card_test

import (
	"testing"
	"time"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/libraries/datetime"
)

func TestCardReminderTime(t *testing.T) {
	timeSource := datetime.NewFakeTimeSource()

	learnCard := card.CreateLearnCardActivity()
	remindCard := card.CreateRemindCardActivity(learnCard)
	today := datetime.GetToday(timeSource)
	timeToRemindExpected := today.Add(time.Minute * 2)
	remindCard.ScheduleTo(timeToRemindExpected)

	timeToRemindActual, err := card.GetTimeToRemindToday(timeSource, remindCard)

	if err != nil {
		t.Fatal(err)
	}

	if timeToRemindActual != timeToRemindExpected {
		t.Fatal("Expected time to remind differs from actual")
	}
}
