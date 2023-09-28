package pull_requests

import (
	"context"

	"github.com/google/go-github/github"
)

type DataAccessLayer interface {
	GatherChecks(string, context.Context) (*github.ListCheckRunsResults, error)
}

type GithubDAL struct {
	client     *github.Client
	owner      string
	repository string
}

func InitGithubDAL(client *github.Client, owner, repository string) DataAccessLayer {

	return &GithubDAL{client: client, owner: owner, repository: repository}
}
