package model

import "database/sql"

type UserAccount struct {
	Id       sql.NullInt64  `json:"id"`
	Email    sql.NullString `json:"email"`
	Password sql.NullString `json:"password"`
	Name     sql.NullString `json:"name"`
	Gender   sql.NullString `json:"gender"`
	AccessID sql.NullInt64  `json:"accessID"`
	Created  sql.NullTime   `json:"created"`
	Modify   sql.NullTime   `json:"modify"`
	IsActive sql.NullBool   `json:"isActive"`
}
