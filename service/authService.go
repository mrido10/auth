package service

import (
	"auth/util"
	"github.com/gin-gonic/gin"
	"log"
)

type authService struct {
	Context *gin.Context
	Error   error
}

func (auth *authService) Initialize(context *gin.Context) *authService {
	auth.Context = context
	return auth
}

func (auth authService) ReadBody(data interface{}) authService {
	if err := auth.Context.ShouldBindJSON(&data); err != nil {
		log.Println(err.Error())
		util.Response(auth.Context, 400, err.Error(), nil)
		auth.Error = err
		return auth
	}
	return auth
}
