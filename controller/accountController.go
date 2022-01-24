package controller

import (
	"auth/dao"
	"auth/model"
	"auth/util"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
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

type Token struct {
	Auth string `json:"authorization"`
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
		if err == sql.ErrNoRows {
			log.Print("Wrong Email or Password!")
			util.Response(c, http.StatusUnauthorized, "Wrong Email or Password!", nil)
			return
		}
		util.Response(c, 400, err.Error(), nil)
		return
	}

	if pwd != acc.Password {
		log.Print("Unauthorized")
		util.Response(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	if acc.IsActive != "YES" {
		fmt.Print("Your account is not active")
		util.Response(c, http.StatusUnauthorized, "Your account is not active", nil)
		return
	}

	token := util.GenerateToken(acc.Name, acc.UserID, acc.AccesID)

	var body Token
	body.Auth = token
	util.Response(c, http.StatusOK, "Succes", body)
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

	err = SendActivationAccount(acc.Email, userID)
	if err != nil {
		fmt.Println(err.Error())
		util.Response(c, 400, err.Error(), nil)
		return
	}

	util.Response(c, 200, "Your account has been registered, please verify your account from your email in 2 minute", nil)

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
