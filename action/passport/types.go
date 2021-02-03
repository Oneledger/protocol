package passport

import "time"

func validateRFCTime(timestamp string) bool {
	timeFormat := time.RFC3339
	_, err := time.Parse(timeFormat, timestamp)
	if err != nil {
		return false
	}
	return true
}
