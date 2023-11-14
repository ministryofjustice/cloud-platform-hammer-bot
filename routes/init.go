package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"go.uber.org/zap"

	ginzap "github.com/gin-contrib/zap"
)

func InitRouter(r *gin.Engine, ghClient *github.Client) {
	InitGetCheckPR(r, ghClient)

	r.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
}

func InitLogger(r *gin.Engine) {
	logger, _ := zap.NewProduction()
	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	r.Use(ginzap.RecoveryWithZap(logger, true))
}
