package init_app

import (
	"github.com/gin-gonic/gin"
	"github.com/ministryofjustice/cloud-platform-hammer-bot/routes"
	"github.com/ministryofjustice/cloud-platform-hammer-bot/utils"
)

func InitGin(gh utils.GitHub) *gin.Engine {

	gin.SetMode(gh.Mode)

	r := gin.New()

	routes.InitLogger(r)

	routes.InitRouter(r, gh)

	return r
}
