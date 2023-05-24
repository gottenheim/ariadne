package datetime

import "time"

func IsDateToday(timeSource TimeSource, date time.Time) bool {
	today := timeSource.Today()
	tomorrow := today.AddDate(0, 0, 1)

	return date.After(today) && date.Before(tomorrow)
}
