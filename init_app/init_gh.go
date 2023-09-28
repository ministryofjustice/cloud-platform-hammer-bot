package init_app

import (
	"context"
	"errors"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func InitGH(token string) (*github.Client, error) {
	var client *github.Client
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client = github.NewClient(tc)

	if token == "" {
		return nil, errors.New("You must provide a github token to avoid rate limits")
	}

	return client, nil
}
