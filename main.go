package main

import (
	"auth/app"
	"auth/util"
	"fmt"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
	"log"
)

func init() {
	logrusInit()
}

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

// logger format
type PlainFormatter struct {
	TimestampFormat string
	LevelDesc []string
}

func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := fmt.Sprintf(entry.Time.Format(f.TimestampFormat))
	level := f.LevelDesc[entry.Level]

	if level != "INFO" && level != "DEBUG"{
		return []byte(fmt.Sprintf("[PROJECT] %s [%s] %s~%v \t %s\n", timestamp, level, entry.Caller.File,
			entry.Caller.Line, entry.Message)), nil
	} else {
		return []byte(fmt.Sprintf("[PROJECT] %s [%s] %s\n", timestamp, level, entry.Message)), nil
	}
}

func logrusInit() {
	plainFormatter := new(PlainFormatter)
	plainFormatter.TimestampFormat = "2006/01/02 - 15:04:05"
	plainFormatter.LevelDesc = []string{"PANIC", "FATAL", "ERROR", "WARN", "INFO", "DEBUG"}
	logrus.SetReportCaller(true)
	logrus.SetFormatter(plainFormatter)
}
