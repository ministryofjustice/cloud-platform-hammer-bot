package init_app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v57/github"
	"github.com/ministryofjustice/cloud-platform-hammer-bot/routes"
)

func InitGin(ginMode string, ghClient *github.Client) *gin.Engine {
	gin.SetMode(ginMode)

	r := gin.New()

	routes.InitLogger(r)

	routes.InitRouter(r, ghClient)

	return r
}
