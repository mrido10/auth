package model

import "database/sql"

type Access struct {
	Id         sql.NullInt64  `json:"id"`
	AccessName sql.NullString `json:"access_name"`
	AccessCode sql.NullString `json:"access_code"`
	Created    sql.NullTime   `json:"created"`
}
