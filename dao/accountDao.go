package dao

import (
	"auth/model"
	"auth/util"
)

func GetUserAccount(email string) (model.UserAccount, error) {
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
		return model.UserAccount{}, err
	}

	return result, nil
}

func InsertUserAccount(acc model.UserAccount) error {
	query := `INSERT INTO userAccount(userID, email, password, name, accesID, isActive, gender)
		VALUE(?, ?, ?, ?, ?, ?, ?)`

	db, err := util.ConnectMySQL()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query(query, acc.UserID, acc.Email, acc.Password, acc.Name, acc.AccesID, acc.IsActive, acc.Gender)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}
