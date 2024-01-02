package pkg

import (
	"time"
)

func IsValidDate(dateString string) bool {
	layout := "2006-01-02"
	location, _ := time.LoadLocation("UTC")

	if _, err := time.ParseInLocation(layout, dateString, location); err != nil {
		return false
	}

	return true
}
