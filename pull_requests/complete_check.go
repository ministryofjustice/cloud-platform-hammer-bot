package pull_requests

import (
	"github.com/google/go-github/github"
)

func CompletedCheck(check *github.CheckRun, prStatus []InvalidChecks) []InvalidChecks {
	switch *check.Conclusion {
	case "success":
		prStatus = nil
	case "skipped":
		prStatus = nil
	case "failure":
		prStatus = append(prStatus, InvalidChecks{*check.Name, "this check failed, check your pr and ammend", Failure, 0})
	case "action_required":
		prStatus = append(prStatus, InvalidChecks{*check.Name, "this check failed because an action is required, check your pr and ammend", Failure, 0})
	case "cancelled":
		prStatus = append(prStatus, InvalidChecks{*check.Name, "this check failed because somebody manually cancelled the check", Failure, 0})
	case "timed_out":
		prStatus = append(prStatus, InvalidChecks{*check.Name, "this check failed because it timed out", Failure, 0})
	case "stale":
		prStatus = append(prStatus, InvalidChecks{*check.Name, "this check failed because it was stale", Failure, 0})
	default:
		prStatus = append(prStatus, InvalidChecks{*check.Name, "unaccounted for state conclusion: " + *check.Conclusion, Failure, 0})
	}
	return prStatus
}
