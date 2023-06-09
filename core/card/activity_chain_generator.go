package card

import (
	"time"

	"github.com/gottenheim/ariadne/libraries/datetime"
)

type GenerateActivity int

const (
	LearnCard                             GenerateActivity = 1
	RemindCard                            GenerateActivity = 2
	CardExecutedYesterday                 GenerateActivity = 4
	CardExecutedMonthAgo                  GenerateActivity = 8
	CardExecutedToday                     GenerateActivity = 16
	RemindCardScheduledToYesterday        GenerateActivity = 32
	RemindCardScheduledToMonthAgo         GenerateActivity = 64
	RemindCardScheduledToTomorrow         GenerateActivity = 128
	RemindCardScheduledToMonthAhead       GenerateActivity = 256
	RemindCardScheduledToToday            GenerateActivity = 512
	RemindCardScheduledToFiveMinutesAgo   GenerateActivity = 1024
	RemindCardScheduledToFiveMinutesAhead GenerateActivity = 2048
)

func GenerateActivityChain(activities ...GenerateActivity) CardActivity {
	var currentActivity CardActivity

	now := datetime.FakeNow()

	for _, activity := range activities {
		if (activity & LearnCard) == LearnCard {
			LearnCard := CreateLearnCardActivity()

			if (activity & CardExecutedMonthAgo) == CardExecutedMonthAgo {
				LearnCard.MarkAsExecuted(now.AddDate(0, -1, 0))
			} else if (activity & CardExecutedYesterday) == CardExecutedYesterday {
				LearnCard.MarkAsExecuted(now.AddDate(0, 0, -1))
			} else if (activity & CardExecutedToday) == CardExecutedToday {
				LearnCard.MarkAsExecuted(now)
			}
			currentActivity = LearnCard
		} else if (activity & RemindCard) == RemindCard {
			RemindCard := CreateRemindCardActivity(currentActivity)

			if (activity & RemindCardScheduledToMonthAgo) == RemindCardScheduledToMonthAgo {
				RemindCard.ScheduleTo(now.AddDate(0, -1, 0))
			} else if (activity & RemindCardScheduledToYesterday) == RemindCardScheduledToYesterday {
				RemindCard.ScheduleTo(now.AddDate(0, 0, -1))
			} else if (activity & RemindCardScheduledToToday) == RemindCardScheduledToToday {
				RemindCard.ScheduleTo(now)
			} else if (activity & RemindCardScheduledToFiveMinutesAgo) == RemindCardScheduledToFiveMinutesAgo {
				RemindCard.ScheduleTo(now.Add(time.Minute * -5))
			} else if (activity & RemindCardScheduledToFiveMinutesAhead) == RemindCardScheduledToFiveMinutesAhead {
				RemindCard.ScheduleTo(now.Add(time.Minute * 5))
			} else if (activity & RemindCardScheduledToTomorrow) == RemindCardScheduledToTomorrow {
				RemindCard.ScheduleTo(now.AddDate(0, 0, 1))
			} else if (activity & RemindCardScheduledToMonthAhead) == RemindCardScheduledToMonthAhead {
				RemindCard.ScheduleTo(now.AddDate(0, 1, 0))
			}

			if (activity & CardExecutedMonthAgo) == CardExecutedMonthAgo {
				RemindCard.MarkAsExecuted(now.AddDate(0, -1, 0))
			} else if (activity & CardExecutedYesterday) == CardExecutedYesterday {
				RemindCard.MarkAsExecuted(now.AddDate(0, 0, -1))
			} else if (activity & CardExecutedToday) == CardExecutedToday {
				RemindCard.MarkAsExecuted(now)
			}
			currentActivity = RemindCard
		}
	}

	return currentActivity
}
