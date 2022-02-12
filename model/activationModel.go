package model

type Activation struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Exp   int64  `json:"exp"`
}