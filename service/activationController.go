package service

import (
	"auth/dao"
	"auth/util"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/gin-gonic/gin"
)

type email struct {
	Email string `json:"email"`
}

func AccountActivate(c *gin.Context) {
	textEncrypted := c.Query("d")

	jsonString, err := util.DecryptData(textEncrypted)
	if err != nil {
		log.Error(err)
		return
	}

	jsonData := []byte(jsonString)

	var data activate

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Error(err.Error())
		return
	}

	intExp := data.Exp
	sec := intExp / 1000
	msec := intExp % 1000
	tExp := time.Unix(sec, msec*int64(time.Millisecond))

	tNow := time.Now()

	if !tNow.Before(tExp) {
		util.Response(c, 400, "Account activation has expired", nil)
		return
	}

	err = dao.UpdateActivateUserAccount(data.ID, data.Email)
	if err != nil {
		log.Error(err.Error())
		util.Response(c, 400, err.Error(), nil)
		return
	}

	util.Response(c, 200, "Your account has been activated", nil)
}

func ReSendActivation(c *gin.Context) {
	var data email
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Error(err)
		util.Response(c, 400, err.Error(), nil)
		return
	}

	acc, err := dao.GetUserAccount(data.Email)
	if err != nil {
		log.Error(err)
		util.Response(c, 400, err.Error(), nil)
		return
	}

	if acc.IsActive.Bool {
		msg := "You account has been activated, you don't need activate again"
		log.Error(msg)
		util.Response(c, 200, msg, nil)
		return
	}

	var aa authService
	err = aa.SendActivationAccount(acc.Email.String, acc.Id.Int64)
	if err != nil {
		log.Error(err.Error())
		util.Response(c, 400, err.Error(), nil)
		return
	}

	util.Response(c, 200, "Resend activation account success, please check your email in 2 minute", nil)
}


