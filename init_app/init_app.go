package init_app

import "github.com/gin-gonic/gin"

func InitGin() *gin.Engine {
	ginMode := "debug"

	gin.SetMode(ginMode)

	r := gin.New()

	// TODO: add a logger

	return r
}
