package pull_requests

import (
	"context"

	"github.com/google/go-github/github"
)

type DataAccessLayer interface {
	ListChecks(string, context.Context) (*github.ListCheckRunsResults, error)
}

type GithubDAL struct {
	client     *github.Client
	owner      string
	repository string
}

// TODO: Might not need this function may delete later
func InitGithubDAL(client *github.Client, owner, repository string) DataAccessLayer {

	return &GithubDAL{client: client, owner: owner, repository: repository}
}
