package controller

import (
	"auth/dao"
	"auth/model"
	"auth/util"
	"bytes"
	"fmt"
	"log"
	"net/http"
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

	linkTo := fmt.Sprintf("http://localhost:3003/activate?userID=%s&email=%s", userID, acc.Email)

	tmpl := "util/html/email-template.html"
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	u := struct{ URL string }{URL: linkTo}

	if err := t.Execute(&tpl, u); err != nil {
		log.Println(err)
		return
	}
	// msg := "Your account has been registered. To activate your account, please click this link "

	subject := "Activate your account"
	err = util.SendEmail(data.Email, "", subject, tpl.String(), tmpl)
	if err != nil {
		fmt.Println(err.Error())
		util.Response(c, 400, err.Error(), nil)
		return
	}

	util.Response(c, 200, "Your account has been registered, please verify your account from your email", nil)

}

func AccountActivate(c *gin.Context) {
	userID := c.Query("userID")
	email := c.Query("email")

	err := dao.UpdateActivateUserAccount(userID, email)
	if err != nil {
		fmt.Println(err.Error())
		util.Response(c, 400, err.Error(), nil)
		return
	}

	util.Response(c, 200, "Your account has been activated", nil)
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
