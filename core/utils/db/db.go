package dbUtil

import (
	"database/sql"
	"github.com/kataras/golog"
	_ "github.com/mattn/go-sqlite3"
	"labTelegramBot/config"
)

func NewDataBase() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", "data/students.db")
	if err != nil {
		golog.Error("Database failed: ", err)
		return nil, err
	}
	golog.Info("Database started")
	err = db.Ping()
	if err != nil {
		golog.Error("Database doesnt replay: ", err)
		return nil, err
	}
	golog.Info("Database available")

	err = initDB(db)
	return db, err
}

func initDB(db *sql.DB) error {
	init := config.GetDBConfig()

	_, err := db.Exec(init)
	return err
}
