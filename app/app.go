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
	route.POST("/login", service.LoginService{}.Login)
	route.POST("/resendActivation", service.ActivationService{}.ReSendActivation)
	route.GET("/activate", service.ActivationService{}.AccountActivate)

	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	if err := route.Run(":" + c.Server.ServicePort); err != nil {
		panic(err)
	}
}
