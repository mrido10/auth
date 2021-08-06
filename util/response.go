package util

import (
	"github.com/gin-gonic/gin"
)

type Resp struct {
	Message string      `json:"message"`
	Status  bool        `json:"status"`
	Data    interface{} `json:"data"`
}

func Response(c *gin.Context, statusCode int, msg string, data interface{}) {
	var resp Resp
	resp.Message = msg
	resp.Data = data
	if statusCode == 200 {
		resp.Status = true
	} else {
		resp.Status = false
	}

	c.JSON(statusCode, resp)
}
