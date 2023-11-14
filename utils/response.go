package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int      `json:"status,omitempty"`
	Message []string `json:"message,omitempty"`
	Error   []string `json:"error,omitempty"`
	Data    any      `json:"data,omitempty"`
}

func SendResponse(c *gin.Context, response Response) {
	var emptyDataObj interface{}
	if len(response.Message) > 0 {
		c.JSON(response.Status, map[string]interface{}{"message": strings.Join(response.Message, "; ")})
		return
	} else if response.Data != emptyDataObj {
		c.JSON(response.Status, response.Data)
		return
	} else if len(response.Error) > 0 {
		c.JSON(response.Status, map[string]interface{}{"error": strings.Join(response.Error, "; ")})
		return
	}
	c.Status(response.Status)
}
