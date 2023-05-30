package datetime

import "time"

func IsToday(timeSource TimeSource, date time.Time) bool {
	today := timeSource.Today()
	tomorrow := today.AddDate(0, 0, 1)

	return date.After(today) && date.Before(tomorrow)
}

func IsBeforeTomorrow(timeSource TimeSource, date time.Time) bool {
	today := timeSource.Today()
	tomorrow := today.AddDate(0, 0, 1)

	return date.Before(tomorrow)
}
