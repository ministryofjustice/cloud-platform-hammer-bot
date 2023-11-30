package routes

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/ministryofjustice/cloud-platform-hammer-bot/utils"

	"github.com/ministryofjustice/cloud-platform-hammer-bot/pull_requests"
)

func InitGetCheckPR(r *gin.Engine, ghClient *github.Client) {
	r.GET("/check-pr/:pr-number", func(c *gin.Context) {
		prNumber := c.Param("pr-number")

		// grab the latest sha for a pr
		s, _ := strconv.Atoi(prNumber)
		pr, _, _ := ghClient.PullRequests.Get(c, "ministryofjustice", "cloud-platform-environments", s)

		status, _, _ := ghClient.Repositories.GetCombinedStatus(c, "ministryofjustice", "cloud-platform-environments", "refs/pull/"+prNumber+"/head", &github.ListOptions{})

		// TODO: map the output of GetCombinedStatus into the InvalidChecks struct, then append to the return data of CheckPRStatus

		log.Println("pr", pr.GetUpdatedAt()) // use this to generate a retryin value
		log.Println("status", status)

		checks, resp, ghErr := ghClient.Checks.ListCheckRunsForRef(c, "ministryofjustice", "cloud-platform-environments", "refs/pull/"+prNumber+"/head", &github.ListCheckRunsOptions{Filter: github.String("all")})
		// checks, resp, ghErr := ghClient.Checks.ListCheckRunsForRef(c, "ministryofjustice", "cloud-platform-environments", *pr.GetHead().SHA, &github.ListCheckRunsOptions{Filter: github.String("all")})

		if ghErr != nil {
			obj := utils.Response{
				Status: resp.StatusCode,
				Error:  []string{ghErr.Error()},
			}
			utils.SendResponse(c, obj)
		}

		data := pull_requests.CheckPRStatus(checks, time.Since)

		obj := utils.Response{
			Status: http.StatusOK,
			Data:   data,
		}
		utils.SendResponse(c, obj)
	})
}
