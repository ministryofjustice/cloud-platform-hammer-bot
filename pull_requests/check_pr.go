package pull_requests

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
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

func CheckCombinedStatus(c *gin.Context, ghClient *github.Client, status *github.CombinedStatus, prNumber string, getTimeSince func(time.Time) time.Duration) []InvalidChecks {
	var statuses []InvalidChecks

	if status.GetState() == "pending" {
		prInt, _ := strconv.Atoi(prNumber)
		pr, _, _ := ghClient.PullRequests.Get(c, "ministryofjustice", "cloud-platform-environments", prInt)
		youngerThan10Mins, timeSinceStart, tenMins := utils.TimeSince(pr.GetUpdatedAt(), getTimeSince)

		if youngerThan10Mins {
			statuses = append(statuses, InvalidChecks{"concourse-ci/status", "this check has been pending for less than 10 minutes, check back again in " + (tenMins - timeSinceStart).String(), Pending, tenMins - timeSinceStart}) // need to calculate the retry in nanoseconds
		} else if !youngerThan10Mins {
			statuses = append(statuses, InvalidChecks{"concourse-ci/status", "this check has been pending for at least 10 minutes, looks like something has gone wrong", Pending, 0})
		}
	}

	for _, s := range status.Statuses {
		if s.GetState() == "failure" {
			statuses = append(statuses, InvalidChecks{s.GetContext(), "this check failed, check your pr and ammend", Failure, 0})
			continue
		}

		if s.GetState() == "success" {
			continue
		}

	}

	return statuses
}
