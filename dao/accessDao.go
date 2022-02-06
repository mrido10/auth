package dao

import (
	"auth/model"
	"auth/util"
	"fmt"
)

func GetAccess(id int64) (result model.Access, err error) {
	query := fmt.Sprintf(`SELECT access_name, access_code FROM access WHERE id = $1`)
	db, err := util.ConnectPostgreSQL()

	if err != nil {
		return result, err
	}
	defer db.Close()

	err = db.QueryRow(query, id).Scan(
		&result.AccessName,
		&result.AccessCode,
	)

	if err != nil {
		return result, err
	}

	return result, nil
}
