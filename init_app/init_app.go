package init_app

import (
	"github.com/gin-gonic/gin"
	"github.com/ministryofjustice/cloud-platform-hammer-bot/routes"
	"github.com/ministryofjustice/cloud-platform-hammer-bot/utils"
)

func InitGin(mode string, gh utils.GitHub) *gin.Engine {
	gin.SetMode(mode)

	r := gin.New()

	routes.InitLogger(r)

	routes.InitRouter(r, gh)

	return r
}
