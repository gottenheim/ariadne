package card_test

import (
	"testing"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/test"
)

type testCardActivity int

const (
	learnCard                       testCardActivity = 1
	remindCard                      testCardActivity = 2
	cardExecutedYesterday           testCardActivity = 4
	cardExecutedMonthAgo            testCardActivity = 8
	cardExecutedToday               testCardActivity = 16
	remindCardScheduledToYesterday  testCardActivity = 32
	remindCardScheduledToMonthAgo   testCardActivity = 64
	remindCardScheduledToTomorrow   testCardActivity = 128
	remindCardScheduledToMonthAhead testCardActivity = 256
	remindCardScheduledToToday      testCardActivity = 512
)

func createTestActivityChain(activities ...testCardActivity) card.CardActivity {
	var currentActivity card.CardActivity

	today := test.GetLocalTestTime()

	for _, activity := range activities {
		if (activity & learnCard) == learnCard {
			learnCard := card.CreateLearnCardActivity()

			if (activity & cardExecutedMonthAgo) == cardExecutedMonthAgo {
				learnCard.MarkAsExecuted(today.AddDate(0, -1, 0))
			} else if (activity & cardExecutedYesterday) == cardExecutedYesterday {
				learnCard.MarkAsExecuted(today.AddDate(0, 0, -1))
			} else if (activity & cardExecutedToday) == cardExecutedToday {
				learnCard.MarkAsExecuted(today)
			}
			currentActivity = learnCard
		} else if (activity & remindCard) == remindCard {
			remindCard := card.CreateRemindCardActivity(currentActivity)

			if (activity & remindCardScheduledToMonthAgo) == remindCardScheduledToMonthAgo {
				remindCard.ScheduleTo(today.AddDate(0, -1, 0))
			} else if (activity & remindCardScheduledToYesterday) == remindCardScheduledToYesterday {
				remindCard.ScheduleTo(today.AddDate(0, 0, -1))
			} else if (activity & remindCardScheduledToToday) == remindCardScheduledToToday {
				remindCard.ScheduleTo(today)
			} else if (activity & remindCardScheduledToTomorrow) == remindCardScheduledToTomorrow {
				remindCard.ScheduleTo(today.AddDate(0, 0, 1))
			} else if (activity & remindCardScheduledToMonthAhead) == remindCardScheduledToMonthAhead {
				remindCard.ScheduleTo(today.AddDate(0, 1, 0))
			}

			if (activity & cardExecutedMonthAgo) == cardExecutedMonthAgo {
				remindCard.MarkAsExecuted(today.AddDate(0, -1, 0))
			} else if (activity & cardExecutedYesterday) == cardExecutedYesterday {
				remindCard.MarkAsExecuted(today.AddDate(0, 0, -1))
			} else if (activity & cardExecutedToday) == cardExecutedToday {
				remindCard.MarkAsExecuted(today)
			}
			currentActivity = remindCard
		}
	}

	return currentActivity
}

func createSerializedTestActivityChain(t *testing.T, activities ...testCardActivity) []byte {
	chain := createTestActivityChain(activities...)

	activitiesBinary, err := card.SerializeCardActivityChain(chain)
	if err != nil {
		t.Fatal(err)
	}

	return activitiesBinary
}
