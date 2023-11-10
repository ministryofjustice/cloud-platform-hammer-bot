package pull_requests

import (
	"time"

	"github.com/google/go-github/github"
)

func InProgressCheck(check *github.CheckRun, prStatus []PRStatus, getTimeSince func(time.Time) time.Duration) []PRStatus {

	timeSince, timeSinceStart, tenMins := timeSince(check, getTimeSince)

	// if the checks are like less 5 mins old, respond to the slack bot to add a comment saying in future make sure to only post your pr when its completed all it's checks
	if timeSince {
		prStatus = append(prStatus, PRStatus{*check.Name, "this check is in_progress and has just been started. check back again in " + (tenMins - timeSinceStart).String(), Pending, tenMins - timeSinceStart})
	} else if !timeSince {
		prStatus = append(prStatus, PRStatus{*check.Name, "this check has been in_progress for at least 10 mins, looks like something has gone wrong?", Pending, 0})
	}
	return prStatus
}

func QueuedCheck(check *github.CheckRun, prStatus []PRStatus, getTimeSince func(time.Time) time.Duration) []PRStatus {

	timeSince, timeSinceStart, tenMins := timeSince(check, getTimeSince)

	// if the checks are like less 5 mins old, respond to the slack bot to add a comment saying in future make sure to only post your pr when its completed all it's checks
	if timeSince {
		prStatus = append(prStatus, PRStatus{*check.Name, "this check is queued and has just started, check back again in " + (tenMins - timeSinceStart).String(), Pending, tenMins - timeSinceStart})
	} else if !timeSince {
		prStatus = append(prStatus, PRStatus{*check.Name, "this check has been queued for at least 10 mins, looks like something has gone wrong?", Pending, 0})
	}
	return prStatus
}

func timeSince(check *github.CheckRun, getTimeSince func(time.Time) time.Duration) (bool, time.Duration, time.Duration) {
	tenMins := time.Duration(10 * time.Minute)
	timeSinceStart := getTimeSince(check.StartedAt.Time)
	rounded := time.Duration(timeSinceStart.Seconds()+0.5) * time.Second
	if (timeSinceStart) < tenMins {
		return true, rounded, tenMins
	} else {
		return false, rounded, tenMins
	}
}
