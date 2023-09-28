package pull_requests

import "github.com/google/go-github/github"

type CheckList struct {
	Name,
	Status,
	Conclusion []string
}

// TODO: maybe we can just delete this
func MapChecks(checks *github.ListCheckRunsResults) CheckList {
	var checklist CheckList
	for _, check := range checks.CheckRuns {
		checklist = CheckList{append(checklist.Name, *check.Name), append(checklist.Status, *check.Status), append(checklist.Conclusion, *check.Conclusion)}
	}
	return checklist
}
