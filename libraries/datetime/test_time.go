package datetime

import "time"

func FakeNow() time.Time {
	return time.Date(2010, 1, 2, 15, 30, 10, 20, time.Local)
}
