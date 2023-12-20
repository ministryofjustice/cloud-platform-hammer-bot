package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/ministryofjustice/cloud-platform-hammer-bot/commit"
	"github.com/ministryofjustice/cloud-platform-hammer-bot/utils"
)

var (
	user = os.Getenv("GITHUB_USER")
	url  = os.Getenv("GITHUB_URL")
)

func InitPostCommit(r *gin.Engine, ghRepo *git.Repository, ghToken string) {
	r.GET("/check-pr", func(c *gin.Context) {
		branch := c.Query("branch")
		repo, err := commit.OpenRepo()
		if err != nil {
			repo, err = commit.CloneRepo(url)
			if err != nil {
				obj := utils.Response{
					Status: 0,
					Error:  []string{err.Error()},
				}
				utils.SendResponse(c, obj)
				return
			}
		}

		err = commit.FetchBranch(repo, branch)
		if err != nil {
			obj := utils.Response{
				Status: 0,
				Error:  []string{err.Error()},
			}
			utils.SendResponse(c, obj)
			return
		}

		err = commit.CheckoutBranch(repo, branch)
		if err != nil {
			obj := utils.Response{
				Status: 0,
				Error:  []string{err.Error()},
			}
			utils.SendResponse(c, obj)
			return
		}

		err = commit.PushCommit(repo, user, ghToken, branch)
		if err != nil {
			obj := utils.Response{
				Status: 0,
				Error:  []string{err.Error()},
			}
			utils.SendResponse(c, obj)
			return
		}
	})
}
