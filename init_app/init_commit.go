package init_app

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/ministryofjustice/cloud-platform-hammer-bot/commit"
)

func InitCommit() (*git.Repository, error) {
	url := os.Getenv("GITHUB_URL")
	repo, err := commit.CloneRepo(url)
	if err != nil {
		panic(err)
	}

	return repo, err
}
