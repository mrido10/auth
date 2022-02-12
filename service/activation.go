package service

import (
	"auth/dao"
	"auth/model"
	"auth/util"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

type ActivationService struct {
	authService
}

var activation model.Activation

func (auth ActivationService) AccountActivate(context *gin.Context) {
	auth.authService.Initialize(context)
	paramData, err := auth.getParameterData()
	if err != nil {
		return
	}

	err = auth.activateUserAccount(paramData, auth.validateDataParameter)
	if err != nil {
		return
	}
	util.Response(auth.authService.Context, 200, "Your account has been activated", nil)
	log.Info("Account has been activated")
}

func (auth ActivationService) getParameterData() (paramData string, err error){
	textEncrypted := auth.authService.Context.Query("d")
	paramData, err = util.DecryptData(textEncrypted)
	if err != nil {
		log.Error(err.Error())
		util.Response(auth.Context, 400, err.Error(), nil)
	}
	return
}

func (auth ActivationService) validateDataParameter(paramData string) (err error){
	err = json.Unmarshal([]byte(paramData), &activation)
	if err != nil {
		log.Error(err.Error())
		util.Response(auth.Context, 400, err.Error(), nil)
	}

	intExp := activation.Exp
	second := intExp / 1000
	miliSecond := intExp % 1000
	tExp := time.Unix(second, miliSecond * int64(time.Millisecond))

	if !time.Now().Before(tExp) {
		err = errors.New("Account activation has expired")
		util.Response(auth.authService.Context, 400, err.Error(), nil)
		return
	}
	return
}

func (auth ActivationService) activateUserAccount(paramData string, validate func(paramData string) error) (err error) {
	if err := validate(paramData); err != nil {
		return err
	}

	err = dao.UpdateActivateUserAccount(activation.ID, activation.Email)
	if err != nil {
		log.Error(err.Error())
		util.Response(auth.authService.Context, 400, err.Error(), nil)
		return
	}

	return
}