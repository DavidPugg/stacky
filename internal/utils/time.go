package utils

import (
	"fmt"
	"time"
)

func FormatCreateTime(timeString string) string {
	t, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return ""
	}

	timeType := "s"
	timeSince := time.Since(t).Round(time.Second).Seconds()

	if timeType == "s" && timeSince >= 60 {
		timeType = "m"
		timeSince = timeSince / 60
	}

	if timeType == "m" && timeSince >= 60 {
		timeType = "h"
		timeSince = timeSince / 60
	}

	if timeType == "h" && timeSince >= 24 {
		timeType = "d"
		timeSince = timeSince / 24
	}

	if timeType == "w" && timeSince >= 7 {
		timeType = "w"
		timeSince = timeSince / 7
	}

	return fmt.Sprintf("%d%s", int(timeSince), timeType)
}
