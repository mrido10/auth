package app

import (
	"auth/config"
	"auth/service"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartService() {
	route := gin.Default()
	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	route.POST("/register", service.RegisterService{}.Register)
	route.POST("/login", service.Login)
	route.POST("/resendActivation", service.ReSendActivation)
	route.GET("/activate", service.AccountActivate)

	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	if err := route.Run(":" + c.Server.ServicePort); err != nil {
		panic(err)
	}
}
