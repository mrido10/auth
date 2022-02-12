package service

import (
	"auth/dao"
	"auth/dto"
	"auth/model"
	"auth/util"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type RegisterService struct {
	authService
}

var register dto.Register
func (auth RegisterService) Register(context *gin.Context) {
	err := auth.authService.Initialize(context).ReadBody(&register).Error
	if err != nil {
		return
	}
	err = auth.checkUserAccount()
	if err != nil {
		return
	}

	err = auth.insertDataRegister()
	if err != nil {
		return
	}

	util.Response(auth.authService.Context, 200,
		"Your account has been registered, please verify your account from your email in 2 minute", nil)
	log.Info("Register Success")
}

func (auth RegisterService) checkUserAccount() (err error) {
	count, err := dao.GetCountUserAccount(register.Email)
	if err != nil {
		log.Error(err.Error())
		util.Response(auth.Context, http.StatusFound, err.Error(), nil)
		return
	}
	if count > 0 {
		err = errors.New("email already exist")
		log.Error(err.Error())
		util.Response(auth.Context, http.StatusFound, err.Error(), nil)
		return
	}
	return
}

func (auth RegisterService) insertDataRegister() (err error){
	password := util.GenerateHmacSHA256(register.Password)
	acc := model.UserAccount{
		Email:    sql.NullString{String: register.Email, Valid: true},
		Password: sql.NullString{String: password, Valid: true},
		Name:     sql.NullString{String: register.Name, Valid: true},
		AccessID: sql.NullInt64{Int64: register.AccessID, Valid: true},
		IsActive: sql.NullBool{Bool: false, Valid: true},
		Gender:   sql.NullString{String: register.Gender, Valid: true},
	}

	err = dao.InsertUserAccount(acc)
	if err != nil {
		log.Error(err.Error())
		util.Response(auth.Context, 400, err.Error(), nil)
		return
	}

	err = auth.authService.SendActivationAccount(acc.Email.String, acc.Id.Int64)
	if err != nil {
		log.Error(err.Error())
		util.Response(auth.Context, 400, err.Error(), nil)
		return
	}
	return
}

