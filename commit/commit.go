package commit

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func CloneRepo(url string) (*git.Repository, error) {
	r, err := git.PlainClone("/app/environments", false, &git.CloneOptions{
		URL:      url,
		Progress: nil,
	})
	if err != nil {
		return nil, fmt.Errorf("an error occurred at clone: %v", err)
	}
	return r, err

}

func OpenRepo() (*git.Repository, error) {
	r, err := git.PlainOpen("/app/environments")
	if err != nil {
		return nil, fmt.Errorf("an error occurred at open: %v", err)
	}
	return r, err
}

func FetchBranch(r *git.Repository, branch string) error {
	ref := fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch)
	err := r.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Progress:   nil,
		RefSpecs:   []config.RefSpec{config.RefSpec(ref)},
	})
	if err != nil {
		return fmt.Errorf("an error occurred at fetch: %v", err)
	}
	return nil
}

func CheckoutBranch(r *git.Repository, branch string) error {
	w, err := r.Worktree()
	if err != nil {
		fmt.Println(err)
	}
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Force:  true,
	})
	if err != nil {
		return fmt.Errorf("an error occurred at checkout: %v", err)
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
		return fmt.Errorf("an error occurred at add: %v", err)
	}

	_, err = w.Commit("Hammer-bot blank commit", &git.CommitOptions{
		AllowEmptyCommits: true,
	})
	if err != nil {
		return fmt.Errorf("an error occurred at commit: %v", err)
	}

	auth := &http.BasicAuth{
		Username: user,
		Password: token,
	}

	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Progress:   nil,
		Auth:       auth,
		ForceWithLease: &git.ForceWithLease{
			RefName: plumbing.NewBranchReferenceName(branch),
		},
	})
	if err != nil {
		return fmt.Errorf("an error occurred at push: %v", err)
	}
	return nil
}
