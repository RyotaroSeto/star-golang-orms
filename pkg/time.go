package pkg

import "time"

func StringToTime(str string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", str)
	return t
}
