package controller

import (
	"auth/config"
	"auth/dao"
	"auth/util"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"text/template"
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
		log.Println(err)
		return
	}

	jsonData := []byte(jsonString)

	var data activate

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println(err.Error())
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
		fmt.Println(err.Error())
		util.Response(c, 400, err.Error(), nil)
		return
	}

	util.Response(c, 200, "Your account has been activated", nil)
}

func ReSendActivation(c *gin.Context) {
	var data email
	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println(err)
		util.Response(c, 400, err.Error(), nil)
		return
	}

	acc, err := dao.GetUserAccount(data.Email)
	if err != nil {
		fmt.Println(err)
		util.Response(c, 400, err.Error(), nil)
		return
	}

	if acc.IsActive.Bool {
		msg := "You account has been activated, you don't need activate again"
		fmt.Println(msg)
		util.Response(c, 200, msg, nil)
		return
	}

	err = SendActivationAccount(acc.Email.String, acc.Id.Int64)
	if err != nil {
		fmt.Println(err.Error())
		util.Response(c, 400, err.Error(), nil)
		return
	}

	util.Response(c, 200, "Resend activation account success, please check your email in 2 minute", nil)
}

func SendActivationAccount(email string, accountID int64) error {
	conf, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
		return err
	}

	exp := time.Now().Add(time.Second*120).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)) // time in milis
	param := fmt.Sprintf(`{"id": "%v", "email":"%v", "exp": %v}`, accountID, email, exp)
	paramEncrypted, err := util.EncryptData(param)
	if err != nil {
		log.Println(err)
		return err
	}

	linkTo := fmt.Sprintf("%s://%s:%s/activate?d=%v", conf.Server.Protocol, conf.Server.Host, conf.Server.ServicePort, url.QueryEscape(paramEncrypted))
	fmt.Println(linkTo)

	tmpl := "util/html/email-template.html"
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	var tpl bytes.Buffer
	u := struct{ URL string }{URL: linkTo}

	if err := t.Execute(&tpl, u); err != nil {
		log.Println(err)
		return err
	}

	subject := "Activate your account"
	err = util.SendEmail(email, "", subject, tpl.String(), tmpl)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
