package pull_requests

import (
	"context"

	"github.com/google/go-github/github"
)

func (c *GithubDAL) ListChecks(reference string, ctx context.Context) (*github.ListCheckRunsResults, error) {
	checks, _, err := c.client.Checks.ListCheckRunsForRef(ctx, c.owner, c.repository, reference, nil)
	if err != nil {
		return nil, err
	}
	return checks, nil
}
