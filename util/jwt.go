package util

import (
	"auth/config"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type authClaims struct {
	Name    string `json:"name"`
	Id      string `json:"id"`
	LoginAs string `json:"loginAs"`
	jwt.StandardClaims
}

func GenerateToken(name string, id string, loginAs string) string {
	claims := &authClaims{
		name,
		id,
		loginAs,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	conf, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	t, err := token.SignedString([]byte(conf.Jwt.Key))
	if err != nil {
		panic(err)
	}
	return t
}
