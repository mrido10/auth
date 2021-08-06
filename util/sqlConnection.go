package util

import (
	"database/sql"
	"fmt"
	"log"

	"auth/config"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectMySQL() (*sql.DB, error) {
	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.MySql.User, c.MySql.Password, c.Server.Host, c.Server.Port, c.MySql.DbName)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
