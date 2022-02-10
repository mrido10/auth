package service

import (
	"auth/dao"
	"auth/util"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type activate struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Exp   int64  `json:"exp"`
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

	if pwd != acc.Password.String {
		log.Print("Unauthorized")
		util.Response(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	if !acc.IsActive.Bool {
		fmt.Print("Your account is not active")
		util.Response(c, http.StatusUnauthorized, "Your account is not active", nil)
		return
	}

	access, err := dao.GetAccess(acc.AccessID.Int64)
	if err != nil {
		util.Response(c, 400, err.Error(), nil)
		return
	}
	token := util.GenerateToken(acc.Name.String, acc.Id.Int64, access.AccessCode.String)

	var body Token
	body.Auth = token
	util.Response(c, http.StatusOK, "Succes", body)
}
