package routes

import (
	"fmt"
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
		checks, resp, ghErr := ghClient.Checks.ListCheckRunsForRef(c, "ministryofjustice", "cloud-platform-environments", "refs/pull/"+prNumber+"/head", &github.ListCheckRunsOptions{Filter: github.String("all")})
		if ghErr != nil {
			obj := utils.Response{
				Status: resp.StatusCode,
				Error:  []string{ghErr.Error()},
			}
			utils.SendResponse(c, obj)
		}

		data, err := pull_requests.CheckPRStatus(checks, time.Since)

		fmt.Printf("checkInvalidChecks %v", data)
		if err != nil {
			obj := utils.Response{
				Status: http.StatusInternalServerError,
				Error:  []string{"Reading from Redis"},
			}
			utils.SendResponse(c, obj)
			return
		}

		obj := utils.Response{
			Status: http.StatusOK,
			Data:   data,
		}
		utils.SendResponse(c, obj)
	})
}
