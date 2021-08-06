package util

import (
	"auth/config"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func GenerateHmacSHA256(data string) string {
	c, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	h := hmac.New(sha256.New, []byte(c.Hash.Secret))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
