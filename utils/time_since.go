package utils

import "time"

func TimeSince(startedAt time.Time, getTimeSince func(time.Time) time.Duration) (bool, time.Duration, time.Duration) {
	tenMins := time.Duration(10 * time.Minute)
	timeSinceStart := getTimeSince(startedAt)
	rounded := time.Duration(timeSinceStart.Seconds()+0.5) * time.Second

	if timeSinceStart < tenMins {
		return true, rounded, tenMins
	} else {
		return false, rounded, tenMins
	}
}
