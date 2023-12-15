package pull_requests

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v57/github"
	"github.com/ministryofjustice/cloud-platform-hammer-bot/utils"
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

func CheckPendingStatus(c *gin.Context, ghClient *github.Client, prNumber string, getTimeSince func(time.Time) time.Duration) (func() InvalidChecks, *github.Response, error) {
	prInt, _ := strconv.Atoi(prNumber)
	pr, resp, ghErr := ghClient.PullRequests.Get(c, "ministryofjustice", "cloud-platform-environments", prInt)
	if ghErr != nil {
		return nil, resp, ghErr
	}

	return func() InvalidChecks {
		youngerThan10Mins, timeSinceStart, tenMins := utils.TimeSince(pr.GetUpdatedAt().Time, getTimeSince)

		if youngerThan10Mins {
			return InvalidChecks{"concourse-ci/status", "this check has been pending for less than 10 minutes, check back again in " + (tenMins - timeSinceStart).String(), Pending, tenMins - timeSinceStart} // need to calculate the retry in nanoseconds
		} else {
			return InvalidChecks{"concourse-ci/status", "this check has been pending for at least 10 minutes, looks like something has gone wrong", Pending, 0}
		}
	}, nil, nil
}

func CheckCombinedStatus(status *github.CombinedStatus, checkPendingFn func() InvalidChecks) []InvalidChecks {
	var statuses []InvalidChecks

	if status.GetState() == "pending" {
		pendingStatus := checkPendingFn()
		statuses = append(statuses, pendingStatus)
	}

	for _, s := range status.Statuses {
		if s.GetState() == "failure" {
			statuses = append(statuses, InvalidChecks{s.GetContext(), "this check failed, check your pr and amend", Failure, 0})
			continue
		}

		if s.GetState() == "success" {
			continue
		}
	}

	return statuses
}
