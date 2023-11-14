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

type InvalidChecks struct {
	Name    string
	Message string
	Status  Status
	RetryIn time.Duration
}

func CheckPRStatus(checks *github.ListCheckRunsResults, getTimeSince func(time.Time) time.Duration) ([]InvalidChecks, error) {
	var prStatus []InvalidChecks
	fmt.Printf("checkruns array %v", checks.GetTotal())

	for _, check := range checks.CheckRuns {
		fmt.Printf("prStatus array %v", prStatus)
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

	return prStatus, nil
}
