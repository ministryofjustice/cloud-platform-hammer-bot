package routes

import (
	// "net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v57/github"
	// "github.com/ministryofjustice/cloud-platform-hammer-bot/utils"
)

func InitGetRetriggerChecks(r *gin.Engine, ghClient *github.Client) {
	r.GET("/retrigger/:pr-number", func(c *gin.Context) {
		// prNumber := c.Param("pr-number")
		// TODO: from the pr number get the github branch
		// TODO: then push an empty commit to that branch

		// TODO: the slackbot waits and then resends the PR -- this is a job for the slack bot

		// obj := utils.Response{
		// 	Status: http.StatusOK,
		// 	Data:   data,
		// }
		// utils.SendResponse(c, obj)
	})
}
