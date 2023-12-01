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

		statuses, _, _ := ghClient.Repositories.GetCombinedStatus(c, "ministryofjustice", "cloud-platform-environments", "refs/pull/"+prNumber+"/head", &github.ListOptions{})

		checks, resp, ghErr := ghClient.Checks.ListCheckRunsForRef(c, "ministryofjustice", "cloud-platform-environments", "refs/pull/"+prNumber+"/head", &github.ListCheckRunsOptions{Filter: github.String("all")})
		// checks, resp, ghErr := ghClient.Checks.ListCheckRunsForRef(c, "ministryofjustice", "cloud-platform-environments", *pr.GetHead().SHA, &github.ListCheckRunsOptions{Filter: github.String("all")})

		if ghErr != nil {
			obj := utils.Response{
				Status: resp.StatusCode,
				Error:  []string{ghErr.Error()},
			}
			utils.SendResponse(c, obj)
		}

		combinedStatus := pull_requests.CheckCombinedStatus(c, ghClient, statuses, prNumber, time.Since)
		data := pull_requests.CheckPRStatus(checks, time.Since)

		data = append(data, combinedStatus...)

		obj := utils.Response{
			Status: http.StatusOK,
			Data:   data,
		}
		utils.SendResponse(c, obj)
	})
}
