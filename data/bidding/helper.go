package bidding

import (
	"time"
)

func timeToString(inputTime time.Time) string {
	str := inputTime.Format(time.RFC3339)
	return str
}
