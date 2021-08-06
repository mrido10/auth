package model

type UserAccount struct {
	UserID   string `json:"userID"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	AccesID  string `json:"accesID"`
	IsActive string `json:"isActive"`
}
