package test

import "time"

func GetLocalTestTime() time.Time {
	return time.Date(2010, 1, 2, 15, 30, 10, 20, time.Local)
}
