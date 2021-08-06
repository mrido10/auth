package app

import (
	"auth/controller"

	"github.com/gin-gonic/gin"
)

func StartService() {
	route := gin.Default()
	route.POST("/login", controller.Login)

	if err := route.Run(":3003"); err != nil {
		panic(err)
	}
}
