package card

import (
	"github.com/gottenheim/ariadne/details/datetime"
)

type GenerateActivity int

const (
	LearnCard                       GenerateActivity = 1
	RemindCard                      GenerateActivity = 2
	CardExecutedYesterday           GenerateActivity = 4
	CardExecutedMonthAgo            GenerateActivity = 8
	CardExecutedToday               GenerateActivity = 16
	RemindCardScheduledToYesterday  GenerateActivity = 32
	RemindCardScheduledToMonthAgo   GenerateActivity = 64
	RemindCardScheduledToTomorrow   GenerateActivity = 128
	RemindCardScheduledToMonthAhead GenerateActivity = 256
	RemindCardScheduledToToday      GenerateActivity = 512
)

func GenerateActivityChain(activities ...GenerateActivity) CardActivity {
	var currentActivity CardActivity

	today := datetime.FakeNow()

	for _, activity := range activities {
		if (activity & LearnCard) == LearnCard {
			LearnCard := CreateLearnCardActivity()

			if (activity & CardExecutedMonthAgo) == CardExecutedMonthAgo {
				LearnCard.MarkAsExecuted(today.AddDate(0, -1, 0))
			} else if (activity & CardExecutedYesterday) == CardExecutedYesterday {
				LearnCard.MarkAsExecuted(today.AddDate(0, 0, -1))
			} else if (activity & CardExecutedToday) == CardExecutedToday {
				LearnCard.MarkAsExecuted(today)
			}
			currentActivity = LearnCard
		} else if (activity & RemindCard) == RemindCard {
			RemindCard := CreateRemindCardActivity(currentActivity)

			if (activity & RemindCardScheduledToMonthAgo) == RemindCardScheduledToMonthAgo {
				RemindCard.ScheduleTo(today.AddDate(0, -1, 0))
			} else if (activity & RemindCardScheduledToYesterday) == RemindCardScheduledToYesterday {
				RemindCard.ScheduleTo(today.AddDate(0, 0, -1))
			} else if (activity & RemindCardScheduledToToday) == RemindCardScheduledToToday {
				RemindCard.ScheduleTo(today)
			} else if (activity & RemindCardScheduledToTomorrow) == RemindCardScheduledToTomorrow {
				RemindCard.ScheduleTo(today.AddDate(0, 0, 1))
			} else if (activity & RemindCardScheduledToMonthAhead) == RemindCardScheduledToMonthAhead {
				RemindCard.ScheduleTo(today.AddDate(0, 1, 0))
			}

			if (activity & CardExecutedMonthAgo) == CardExecutedMonthAgo {
				RemindCard.MarkAsExecuted(today.AddDate(0, -1, 0))
			} else if (activity & CardExecutedYesterday) == CardExecutedYesterday {
				RemindCard.MarkAsExecuted(today.AddDate(0, 0, -1))
			} else if (activity & CardExecutedToday) == CardExecutedToday {
				RemindCard.MarkAsExecuted(today)
			}
			currentActivity = RemindCard
		}
	}

	return currentActivity
}
