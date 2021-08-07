package app

import (
	"auth/controller"

	"github.com/gin-gonic/gin"
)

func StartService() {
	route := gin.Default()
	route.POST("/register", controller.Register)
	route.POST("/login", controller.Login)
	route.GET("/activate", controller.AccountActivate)

	if err := route.Run(":3003"); err != nil {
		panic(err)
	}
}
