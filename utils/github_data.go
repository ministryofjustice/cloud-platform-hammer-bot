package utils

import (
	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v57/github"
)

type GitHub struct {
	Mode,
	Token,
	URL,
	User string
	Repo   *git.Repository
	Client *github.Client
}
