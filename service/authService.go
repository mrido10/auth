package service

import (
	"auth/config"
	"auth/util"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/url"
	"time"
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

func (auth authService) SendActivationAccount(email string, accountID int64) error {
	conf, err := config.GetConfig()
	if err != nil {
		log.Error(err)
		return err
	}

	exp := time.Now().Add(time.Second*120).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)) // time in milis
	param := fmt.Sprintf(`{"id": %v, "email":"%v", "exp": %v}`, accountID, email, exp)
	paramEncrypted, err := util.EncryptData(param)
	if err != nil {
		log.Error(err)
		return err
	}

	linkTo := fmt.Sprintf("%s://%s:%s/activate?d=%v", conf.Server.Protocol, conf.Server.Host, conf.Server.ServicePort, url.QueryEscape(paramEncrypted))
	tmpl := "util/html/email-template.html"
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Error(err)
		return err
	}

	var tpl bytes.Buffer
	u := struct{ URL string }{URL: linkTo}

	if err := t.Execute(&tpl, u); err != nil {
		log.Error(err)
		return err
	}

	subject := "Activate your account"
	err = util.SendEmail(email, "", subject, tpl.String(), tmpl)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
