package service

import (
	"auth/dao"
	"auth/util"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (auth ActivationService) ReSendActivation(context *gin.Context) {
	err := auth.authService.Initialize(context).ReadBody(&activation).Error
	if err != nil {
		return
	}
	if auth.resendActivation() != nil {
		return
	}
	util.Response(auth.authService.Context, 200, "Resend activation account success, please check your email in 2 minute", nil)
	log.Info("Resend activation account success")
}

func (auth ActivationService) resendActivation() (err error){
	account, err := dao.GetUserAccount(activation.Email)
	if err != nil {
		log.Error(err)
		util.Response(auth.authService.Context, 400, err.Error(), nil)
		return
	}

	if account.IsActive.Bool {
		msg := "You account has been activated, you don't need activate again"
		log.Error(msg)
		util.Response(auth.authService.Context, 200, msg, nil)
		err = errors.New("")
		return
	}
	err = auth.authService.SendActivationAccount(account.Email.String, account.Id.Int64)
	if err != nil {
		log.Error(err.Error())
		util.Response(auth.authService.Context, 400, err.Error(), nil)
		return
	}
	return
}
