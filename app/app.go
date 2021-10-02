package app

import (
	"auth/config"
	"auth/controller"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartService() {
	route := gin.Default()
	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	}))

	route.POST("/register", controller.Register)
	route.POST("/login", controller.Login)
	route.POST("/resendActivation", controller.ReSendActivation)
	route.GET("/activate", controller.AccountActivate)

	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	if err := route.Run(":" + c.Server.ServicePort); err != nil {
		panic(err)
	}
}
