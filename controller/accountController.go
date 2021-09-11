package controller

import (
	"auth/config"
	"auth/dao"
	"auth/model"
	"auth/util"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type regist struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RePassword string `json:"rePassword"`
	Name       string `json:"name"`
	Gender     string `json:"gender"`
}

type activate struct {
	UserID string `json:"userID"`
	Email  string `json:"email"`
	Exp    int64  `json:"exp"`
}

func Login(c *gin.Context) {
	var data login
	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println(err)
		util.Response(c, 400, err.Error(), nil)
		return
	}

	pwd := util.GenerateHmacSHA256(data.Password)
	acc, err := dao.GetUserAccount(data.Email)

	if err != nil {
		fmt.Println(err.Error())
		util.Response(c, 400, err.Error(), nil)
		return
	}

	if pwd != acc.Password {
		log.Print("unauthorized")
		util.Response(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	if acc.IsActive != "YES" {
		fmt.Print("your account is not active")
		util.Response(c, http.StatusUnauthorized, "your account is not active", nil)
		return
	}

	token := util.GenerateToken(acc.Name, acc.UserID, acc.AccesID)
	c.Writer.Header().Set("authorization", token)
	util.Response(c, http.StatusOK, "succes", nil)
}

func Register(c *gin.Context) {
	var data regist
	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println(err)
		util.Response(c, 400, err.Error(), nil)
		return
	}

	pwd := util.GenerateHmacSHA256(data.Password)
	_, err := dao.GetUserAccount(data.Email)

	if err == nil {
		fmt.Println("email already exist")
		util.Response(c, http.StatusFound, "email already exist", nil)
		return
	}

	if err.Error() != "sql: no rows in result set" {
		fmt.Println(err.Error())
		util.Response(c, 400, err.Error(), nil)
		return
	}

	userID := generateUserID()

	acc := model.UserAccount{
		UserID:   userID,
		Email:    data.Email,
		Password: pwd,
		Name:     data.Name,
		AccesID:  "cl001",
		IsActive: "NO",
		Gender:   data.Gender,
	}

	err = dao.InsertUserAccount(acc)

	if err != nil {
		fmt.Println(err.Error())
		util.Response(c, 400, err.Error(), nil)
		return
	}

	err = sendActivationAccount(acc.Email, userID)
	if err != nil {
		fmt.Println(err.Error())
		util.Response(c, 400, err.Error(), nil)
		return
	}

	util.Response(c, 200, "Your account has been registered, please verify your account from your email", nil)

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

	userID := data.UserID
	email := data.Email
	intExp := data.Exp

	sec := intExp / 1000
	msec := intExp % 1000
	tExp := time.Unix(sec, msec*int64(time.Millisecond))

	tNow := time.Now()

	if !tNow.Before(tExp) {
		util.Response(c, 400, "Account activation has expired", nil)
		return
	}

	err = dao.UpdateActivateUserAccount(userID, email)
	if err != nil {
		fmt.Println(err.Error())
		util.Response(c, 400, err.Error(), nil)
		return
	}

	util.Response(c, 200, "Your account has been activated", nil)
}

func sendActivationAccount(email string, userID string) error {
	conf, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
		return err
	}

	exp := time.Now().Add(time.Second*60).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)) // time in milis
	param := fmt.Sprintf(`{"userID": "%v", "email":"%v", "exp": %v}`, userID, email, exp)
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

func generateUserID() string {
	t := time.Now()
	tString := t.Format(("2006-01-02 15:04:05.000"))
	tString = strings.Replace(tString, "-", "", 2)
	tString = strings.Replace(tString, " ", "", 1)
	tString = strings.Replace(tString, ":", "", 2)
	tString = strings.Replace(tString, ".", "", 1)
	tString = tString[2 : len(tString)-1]

	return "MM" + tString
}
