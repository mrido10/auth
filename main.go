package main

import (
	"auth/app"
	"auth/util"
	"fmt"
	migrate "github.com/rubenv/sql-migrate"
	"log"
)

func main() {
	runSqlMigrationUP()
	app.StartService()
}

func runSqlMigrationUP() {
	migrations := &migrate.FileMigrationSource{Dir: "sqlMigrations/"}

	db, err := util.ConnectPostgreSQL()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatal("Error occcured:", err)
		return
	}

	log.Println(fmt.Sprintf("Applied %d migrations!\n", n))
}
