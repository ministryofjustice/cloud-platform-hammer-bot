package pull_requests

import (
	"fmt"
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

func CheckPRStatus(checks *github.ListCheckRunsResults) []PRStatus {
	var prStatus []PRStatus // failed or pending statuses

	for _, check := range checks.CheckRuns {
		if *check.Status == "completed" {
			switch *check.Conclusion {
			case "success":
				continue
			case "skipped":
				continue
			case "failure":
				prStatus = append(prStatus, PRStatus{*check.Name, "this check failed check your pr and ammend", Failure, 0})
				continue
			case "action_required":
				prStatus = append(prStatus, PRStatus{*check.Name, "this check failed because an action is required check your pr and ammend", Failure, 0})
				continue
			case "cancelled":
				prStatus = append(prStatus, PRStatus{*check.Name, "this check failed because somebody manually cancelled the check", Failure, 0})
				continue
			default:
				fmt.Printf("unaccounted for value %s %s", *check.Name, *check.Status) // neutral / null
			}
		}

		if *check.Status == "in_progress" {
			// 10 mins - elaspsed time = retry again in x time
			// ((time now - started at time) elaspsed time)
			tenMins := time.Duration(-10 * time.Minute)

			// if the checks are like less 5 mins old, respond to the slack bot to add a comment saying in future make sure to only post your pr when its completed all it's checks
			if (time.Since(*&check.StartedAt.Time)) < tenMins {
				prStatus = append(prStatus, PRStatus{*check.Name, "this check has only just been started check back again in x mins", Pending, time.Since(*&check.StartedAt.Time) - tenMins})
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
//        "neutral",
//        "cancelled",
//        "skipped", x
//        "timed_out", x
//        "action_required", x
//        null
//      ],
