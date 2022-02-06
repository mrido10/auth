package dao

import (
	"auth/model"
	"auth/util"
	"fmt"
)

func GetUserAccount(email string) (model.UserAccount, error) {
	query := fmt.Sprintf(`SELECT id, email, password, is_active FROM user_account WHERE email = $1`)
	db, err := util.ConnectPostgreSQL()
	if err != nil {
		return model.UserAccount{}, err
	}
	defer db.Close()

	var result model.UserAccount
	err = db.QueryRow(query, email).Scan(&result.Id, &result.Email, &result.Password, &result.IsActive)

	if err != nil {
		return model.UserAccount{}, err
	}

	return result, nil
}

func InsertUserAccount(acc model.UserAccount) error {
	query := fmt.Sprintf(`INSERT INTO user_account (email, password, name, gender, acces_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)`)

	db, err := util.ConnectPostgreSQL()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query(query, acc.Email.String, acc.Password.String, acc.Name.String,
		acc.Gender.String, acc.AccessID.Int64, acc.IsActive.Bool)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func UpdateActivateUserAccount(id int64, email string) error {
	query := fmt.Sprintf(`UPDATE user_account SET is_active = true WHERE id = $1 AND email = $2`)

	db, err := util.ConnectPostgreSQL()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query(query, id, email)

	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}
