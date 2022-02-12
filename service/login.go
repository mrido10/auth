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

type LoginService struct {
	authService
}

var login dto.Login
func (auth LoginService) Login(context *gin.Context) {
	err := auth.authService.Initialize(context).ReadBody(&login).Error
	if err != nil {
		return
	}
	token, err := auth.generateToken(auth.validateUserAccount)
	if err != nil {
		return
	}

	util.Response(auth.authService.Context, http.StatusOK, "Success", dto.Token{
		Authorization: token,
	})
	log.Info("Login success")
}

func (auth LoginService) validateUserAccount() (account model.UserAccount, err error){
	account, err = dao.GetUserAccount(login.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("Wrong Email or Password!")
			log.Error(err.Error())
			util.Response(auth.authService.Context, http.StatusUnauthorized, err.Error(), nil)
			return
		}
		log.Error(err.Error())
		util.Response(auth.authService.Context, 400, err.Error(), nil)
		return
	}

	password := util.GenerateHmacSHA256(login.Password)
	if password != account.Password.String {
		err = errors.New("Wrong Email or Password!")
		log.Error(err.Error())
		util.Response(auth.authService.Context, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	if !account.IsActive.Bool {
		err = errors.New("Your account is not active")
		log.Error(err.Error())
		util.Response(auth.authService.Context, http.StatusUnauthorized, err.Error(), nil)
		return
	}
	return
}

func (auth LoginService) generateToken(validate func() (model.UserAccount, error)) (token string, err error) {
	account, err := validate()
	if err != nil {
		return
	}

	access, err := dao.GetAccess(account.AccessID.Int64)
	if err != nil {
		log.Error(err.Error())
		util.Response(auth.authService.Context, 400, err.Error(), nil)
		return
	}
	token = util.GenerateToken(account.Name.String, account.Id.Int64, access.AccessCode.String)
	return
}
