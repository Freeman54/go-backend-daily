package retryafterhint

import (
	"strconv"
	"time"
)

func Seconds(now time.Time, retryAt time.Time, maxSeconds int) string {
	if maxSeconds <= 0 {
		maxSeconds = 1
	}

	wait := retryAt.Sub(now)
	seconds := int((wait + time.Second - 1) / time.Second)
	if seconds < 1 {
		seconds = 1
	}
	if seconds > maxSeconds {
		seconds = maxSeconds
	}
	return strconv.Itoa(seconds)
}
