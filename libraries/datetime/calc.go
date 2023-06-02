package datetime

import "time"

func GetToday(timeSource TimeSource) time.Time {
	now := timeSource.Now()
	return GetDayForDateTime(now)
}

func IsToday(timeSource TimeSource, date time.Time) bool {
	today := GetToday(timeSource)
	tomorrow := today.AddDate(0, 0, 1)

	return date.After(today) && date.Before(tomorrow)
}

func IsSameDay(timeSource TimeSource, leftDate time.Time, rightDate time.Time) bool {
	leftDay := GetDayForDateTime(leftDate)
	rightDay := GetDayForDateTime(rightDate)

	return leftDay == rightDay
}

func IsBeforeTomorrow(timeSource TimeSource, date time.Time) bool {
	today := GetToday(timeSource)
	tomorrow := today.AddDate(0, 0, 1)

	return date.Before(tomorrow)
}

func GetDayForDateTime(dateTime time.Time) time.Time {
	return time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 0, 0, 0, 0, time.Local)
}
