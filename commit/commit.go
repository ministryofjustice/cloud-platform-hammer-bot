package commit

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func CloneRepo(url, user, token string) (*git.Repository, error) {
	auth := &http.BasicAuth{
		Username: user,
		Password: token,
	}

	r, err := git.PlainClone("/app/environments", false, &git.CloneOptions{
		URL:      url,
		Auth:     auth,
		Progress: nil,
	})
	if err != nil {
		return nil, fmt.Errorf("an error occurred at clone: %v", err.Error())
	}
	return r, nil
}

func OpenRepo() (*git.Repository, error) {
	r, err := git.PlainOpen("/app/environments")
	if err != nil {
		return nil, fmt.Errorf("an error occurred at open: %v", err.Error())
	}
	return r, nil
}

func FetchBranch(r *git.Repository, branch string) error {
	ref := fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch)
	err := r.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Progress:   nil,
		RefSpecs:   []config.RefSpec{config.RefSpec(ref)},
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("an error occurred at fetch: %v", err.Error())
	}
	return nil
}

func CheckoutBranch(r *git.Repository, branch string) error {
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("an error occurred at worktree: %v", err.Error())
	}
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Force:  true,
	})
	if err != nil {
		return fmt.Errorf("an error occurred at checkout: %v", err.Error())
	}
	return nil
}

func PushCommit(r *git.Repository, user, token, branch string) error {
	w, err := r.Worktree()
	if err != nil {
		fmt.Println(err)
	}
	_, err = w.Add(".")
	if err != nil {
		return fmt.Errorf("an error occurred at add: %v", err.Error())
	}

	_, err = w.Commit("Hammer-bot blank commit", &git.CommitOptions{
		AllowEmptyCommits: true,
		Author: &object.Signature{
			Name:  "Hammer-bot",
			Email: "jack.stockley@digital.justice.gov.uk",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("an error occurred at commit: %v", err.Error())
	}

	auth := &http.BasicAuth{
		Username: user,
		Password: token,
	}

	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Progress:   nil,
		Auth:       auth,
	})
	if err != nil {
		return fmt.Errorf("an error occurred at push: %v", err.Error())
	}
	return nil
}

func SwitchToMainBranch(r *git.Repository) error {
	w, err := r.Worktree()
	if err != nil {
		fmt.Println(err)
	}
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("main"),
		Force:  true,
	})
	if err != nil {
		return fmt.Errorf("an error occurred at checkout: %v", err.Error())
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Progress:   nil,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("an error occurred at switch to main branch: %v", err.Error())
	}

	return nil
}
