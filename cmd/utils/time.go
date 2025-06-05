package utils

import "time"

type TimeArray [7]int

func (t TimeArray) ToTime() time.Time {
	return time.Date(t[0], time.Month(t[1]), t[2], t[3], t[4], t[5], t[6], time.UTC)
}
