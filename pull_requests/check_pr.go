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

type InvalidChecks struct {
	Name           string
	Message        string
	Status         Status
	RetryInNanoSec time.Duration
}

func CheckPRStatus(checks *github.ListCheckRunsResults, getTimeSince func(time.Time) time.Duration) []InvalidChecks {
	var prStatus []InvalidChecks

	for _, check := range checks.CheckRuns {

		if *check.Status == "completed" {
			prStatus = CompletedCheck(check, prStatus)
			continue
		}

		if *check.Status == "in_progress" {
			prStatus = InProgressCheck(check, prStatus, getTimeSince)
			continue
		}

		if *check.Status == "queued" {
			prStatus = QueuedCheck(check, prStatus, getTimeSince)
			continue
		}
	}

	return prStatus
}
