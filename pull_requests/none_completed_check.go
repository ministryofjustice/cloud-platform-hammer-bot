package pull_requests

import (
	"time"

	"github.com/google/go-github/v57/github"
	"github.com/ministryofjustice/cloud-platform-hammer-bot/utils"
)

func InProgressCheck(check *github.CheckRun, prStatus []InvalidChecks, getTimeSince func(time.Time) time.Duration) []InvalidChecks {
	youngerThan10Mins, timeSinceStart, tenMins := utils.TimeSince(check.StartedAt.Time, getTimeSince)

	if youngerThan10Mins {
		prStatus = append(prStatus, InvalidChecks{*check.Name, "this check is in_progress and has just been started. check back again in " + (tenMins - timeSinceStart).String(), Pending, tenMins - timeSinceStart})
	} else if !youngerThan10Mins {
		prStatus = append(prStatus, InvalidChecks{*check.Name, "this check has been in_progress for at least 10 mins, looks like something has gone wrong?", Pending, 0})
	}
	return prStatus
}

func QueuedCheck(check *github.CheckRun, prStatus []InvalidChecks, getTimeSince func(time.Time) time.Duration) []InvalidChecks {
	youngerThan10Mins, timeSinceStart, tenMins := utils.TimeSince(check.StartedAt.Time, getTimeSince)

	if youngerThan10Mins {
		prStatus = append(prStatus, InvalidChecks{*check.Name, "this check has been queued for less than 10 minutes, check back again in " + (tenMins - timeSinceStart).String(), Pending, tenMins - timeSinceStart})
	} else if !youngerThan10Mins {
		prStatus = append(prStatus, InvalidChecks{*check.Name, "this check has been queued for at least 10 mins, looks like something has gone wrong?", Pending, 0})
	}
	return prStatus
}
