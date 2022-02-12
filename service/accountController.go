package service



type activate struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Exp   int64  `json:"exp"`
}

type Token struct {
	Auth string `json:"authorization"`
}

