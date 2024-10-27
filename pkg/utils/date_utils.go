package utils

import (
	"fmt"
	"time"
)

func FormatTimeToString(date time.Time) string {
	str := date.Format(time.RFC3339)
	return str
}

func FormatStringToTime(str string) time.Time {
	date, err := time.Parse(time.RFC3339, str)
	if err != nil {
		fmt.Println("Could not parse time:", err)
	}
	return date
}
