package app

import (
	"auth/config"
	"auth/controller"
	"log"

	"github.com/gin-gonic/gin"
)

func StartService() {
	route := gin.Default()
	route.POST("/register", controller.Register)
	route.POST("/login", controller.Login)
	route.GET("/activate", controller.AccountActivate)

	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	if err := route.Run(":" + c.Server.ServicePort); err != nil {
		panic(err)
	}
}
