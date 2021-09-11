package util

import (
	"auth/config"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// encrypt using AES
func EncryptData(plainText string) (string, error) {
	conf, err := config.GetConfig()
	if err != nil {
		return "", err
	}

	text := []byte(plainText)
	key := []byte(conf.EncryptDecrypt.Key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], text)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func DecryptData(textEncrypted string) (string, error) {
	conf, err := config.GetConfig()
	if err != nil {
		return "", err
	}
	key := []byte(conf.EncryptDecrypt.Key)

	cipherText, err := base64.URLEncoding.DecodeString(textEncrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return "", err
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
