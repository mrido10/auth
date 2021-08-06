package controller

import (
	"auth/dao"
	"auth/util"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var data login
	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println(err)
		util.Response(c, 400, err.Error(), nil)
		return
	}

	pwd := util.GenerateHmacSHA256(data.Password)
	acc, err := dao.GetAccount(data.Email, pwd)

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
	util.Response(c, 200, "succes", nil)
}
