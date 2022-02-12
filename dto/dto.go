package dto

type Register struct {
	Login
	RePassword string `json:"rePassword"`
	Name       string `json:"name"`
	Gender     string `json:"gender"`
	AccessID   int64  `json:"accessId"`
}

type Login struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type Token struct {
	Authorization string `json:"authorization"`
}