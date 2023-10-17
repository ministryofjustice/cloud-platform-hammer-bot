package pull_requests

import (
	"time"

	"github.com/google/go-github/github"
)

func NoneCompletedCheck(check *github.CheckRun, prStatus []PRStatus, getTimeSince func(time.Time) time.Duration) []PRStatus {
	// 10 mins - elaspsed time = retry again in x time
	tenMins := time.Duration(10 * time.Minute)
	timeSinceStart := getTimeSince(*&check.StartedAt.Time)
	// if the checks are like less 5 mins old, respond to the slack bot to add a comment saying in future make sure to only post your pr when its completed all it's checks
	if (timeSinceStart) < tenMins {
		prStatus = append(prStatus, PRStatus{*check.Name, "this check has only just been started check back again in " + (tenMins - timeSinceStart).String(), Pending, tenMins - timeSinceStart})
	} else {
		prStatus = append(prStatus, PRStatus{*check.Name, "this check has been running for at least 10 mins, looks like something has gone wrong?", Pending, 0})
	}
	return prStatus
}
