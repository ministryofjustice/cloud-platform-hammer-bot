package routes

import (
	"fmt"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ministryofjustice/cloud-platform-hammer-bot/commit"
	"github.com/ministryofjustice/cloud-platform-hammer-bot/utils"
)

func InitGetRetriggerChecks(r *gin.Engine, gh utils.GitHub) {
	r.GET("/retrigger-checks", func(c *gin.Context) {
		b := c.Query("branch")
		splitBranch := strings.Split(b, ",")

		for _, branch := range splitBranch {
			repo, err := commit.OpenRepo()
			if err != nil {
				repo, err = commit.CloneRepo(gh.URL, gh.User, gh.Token)
				if err != nil {
					obj := utils.Response{
						Status: http.StatusInternalServerError,
						Error:  []string{err.Error()},
					}
					utils.SendResponse(c, obj)
					return
				}
			}

			err = commit.FetchBranch(repo, branch)
			if err != nil {
				obj := utils.Response{
					Status: http.StatusInternalServerError,
					Error:  []string{err.Error()},
				}
				utils.SendResponse(c, obj)
				return
			}

			err = commit.CheckoutBranch(repo, branch)
			if err != nil {
				obj := utils.Response{
					Status: http.StatusInternalServerError,
					Error:  []string{err.Error()},
				}
				utils.SendResponse(c, obj)
				return
			}

			err = commit.PushCommit(repo, gh.User, gh.Token, branch)
			if err != nil {
				obj := utils.Response{
					Status: http.StatusInternalServerError,
					Error:  []string{err.Error()},
				}
				utils.SendResponse(c, obj)
				return
			}

			err = commit.SwitchToMainBranch(repo)
			if err != nil {
				obj := utils.Response{
					Status: http.StatusInternalServerError,
					Error:  []string{err.Error()},
				}
				utils.SendResponse(c, obj)
				return
			} else {
				err := fmt.Errorf("retriggered checks for %s completed", branch)
				obj := utils.Response{
					Status: 0,
					Error:  []string{fmt.Sprintf("%v", err)},
				}
				utils.SendResponse(c, obj)
				return
			}
		}
	})
}
