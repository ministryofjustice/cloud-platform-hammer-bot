package pull_requests

import (
	"time"

	"github.com/google/go-github/github"
)

type Status int

const (
	Success Status = iota
	Failure
	Pending
)

type PRStatus struct {
	Name    string
	Message string
	Status  Status
	retryIn time.Duration
}

func CheckPRStatus(checks *github.ListCheckRunsResults, getTimeSince func(time.Time) time.Duration) []PRStatus {
	var prStatus []PRStatus // failed or pending statuses

	for _, check := range checks.CheckRuns {
		if *check.Status == "completed" {
			switch *check.Conclusion {
			case "success":
				continue
			case "skipped":
				continue
			case "failure":
				prStatus = append(prStatus, PRStatus{*check.Name, "this check failed, check your pr and ammend", Failure, 0})
				continue
			case "action_required":
				prStatus = append(prStatus, PRStatus{*check.Name, "this check failed because an action is required, check your pr and ammend", Failure, 0})
				continue
			case "cancelled":
				prStatus = append(prStatus, PRStatus{*check.Name, "this check failed because somebody manually cancelled the check", Failure, 0})
				continue
			case "timed_out":
				prStatus = append(prStatus, PRStatus{*check.Name, "this check failed because it timed out", Failure, 0})
				continue
			case "stale":
				prStatus = append(prStatus, PRStatus{*check.Name, "this check failed because it was stale", Failure, 0})
				continue
			default:
				prStatus = append(prStatus, PRStatus{*check.Name, "unaccounted for state conclusion: " + *check.Conclusion, Failure, 0})
				continue
			}
		}

		if *check.Status == "in_progress" {
			// 10 mins - elaspsed time = retry again in x time
			tenMins := time.Duration(10 * time.Minute)
			timeSinceStart := getTimeSince(*&check.StartedAt.Time)

			// if the checks are like less 5 mins old, respond to the slack bot to add a comment saying in future make sure to only post your pr when its completed all it's checks
			if (timeSinceStart) < tenMins {
				prStatus = append(prStatus, PRStatus{*check.Name, "this check has only just been started check back again in " + (tenMins - timeSinceStart).String(), Pending, tenMins - timeSinceStart})
				continue
			} else {
				prStatus = append(prStatus, PRStatus{*check.Name, "this check has been running for at least 10 mins, looks like something has gone wrong?", Pending, 0})
				continue
			}
		}

		// TODO:  handle queued checks (this is important for concourse and knowing when concourse has stalled)
	}

	return prStatus
}

// "type": "string",
//      "enum": [
//        "queued",
//        "in_progress", x
//        "completed" x
//      ],
//      "examples": [
//        "queued"
//      ]
//    },
//    "conclusion": {
//      "type": [
//        "string",
//        "null"
//      ],
//      "enum": [
//        "success", x
//        "failure", x
//        "neutral", x
//        "cancelled",
//        "skipped", x
//        "stale", x
//        "timed_out", x
//        "action_required", x
//        null x
//      ],
