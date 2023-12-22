package init_app

import (
	"github.com/go-git/go-git/v5"
	"github.com/ministryofjustice/cloud-platform-hammer-bot/commit"
)

func InitCommit(url string) (*git.Repository, error) {
	repo, err := commit.CloneRepo(url)
	if err != nil {
		panic(err)
	}

	return repo, err
}
