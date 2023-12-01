package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/ministryofjustice/cloud-platform-hammer-bot/utils"

	"github.com/ministryofjustice/cloud-platform-hammer-bot/pull_requests"
)

func InitGetCheckPR(r *gin.Engine, ghClient *github.Client) {
	r.GET("/check-pr/:pr-number", func(c *gin.Context) {
		prNumber := c.Param("pr-number")

		statuses, statusResp, ghStatusErr := ghClient.Repositories.GetCombinedStatus(c, "ministryofjustice", "cloud-platform-environments", "refs/pull/"+prNumber+"/head", &github.ListOptions{})
		if ghStatusErr != nil {
			obj := utils.Response{
				Status: statusResp.StatusCode,
				Error:  []string{ghStatusErr.Error()},
			}
			utils.SendResponse(c, obj)
			return
		}

		checks, resp, ghErr := ghClient.Checks.ListCheckRunsForRef(c, "ministryofjustice", "cloud-platform-environments", "refs/pull/"+prNumber+"/head", &github.ListCheckRunsOptions{Filter: github.String("all")})

		if ghErr != nil {
			obj := utils.Response{
				Status: resp.StatusCode,
				Error:  []string{ghErr.Error()},
			}
			utils.SendResponse(c, obj)
			return
		}

		pendingStatusFn, prResp, ghPRErr := pull_requests.CheckPendingStatus(c, ghClient, prNumber, time.Since)
		if ghPRErr != nil {
			obj := utils.Response{
				Status: prResp.StatusCode,
				Error:  []string{ghPRErr.Error()},
			}
			utils.SendResponse(c, obj)
			return
		}

		combinedStatus := pull_requests.CheckCombinedStatus(statuses, pendingStatusFn)
		data := pull_requests.CheckPRStatus(checks, time.Since)

		data = append(data, combinedStatus...)

		obj := utils.Response{
			Status: http.StatusOK,
			Data:   data,
		}
		utils.SendResponse(c, obj)
		return
	})
}
