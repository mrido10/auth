package dao

import (
	"auth/model"
	"auth/util"
	"fmt"
)

func GetAccount(email string, password string) (model.UserAccount, error) {
	query := `SELECT userID, email, password, name, accesID, isActive FROM userAccount
		WHERE email = ?`

	db, err := util.ConnectMySQL()
	if err != nil {
		return model.UserAccount{}, err
	}
	defer db.Close()

	var result model.UserAccount
	err = db.QueryRow(query, email).Scan(
		&result.UserID,
		&result.Email,
		&result.Password,
		&result.Name,
		&result.AccesID,
		&result.IsActive,
	)

	if err != nil {
		fmt.Println(err.Error())
		return model.UserAccount{}, err
	}

	return result, nil
}
