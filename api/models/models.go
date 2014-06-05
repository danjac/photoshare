package models

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
)

var dbMap *gorp.DbMap

func populateDatabase() error {

	dbMap.AddTableWithName(User{}, "users").SetKeys(true, "ID")
	dbMap.AddTableWithName(Photo{}, "photos").SetKeys(true, "ID")
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		panic(err)
	}

	numUsers, err := dbMap.SelectInt("SELECT COUNT(id) FROM users")
	if err != nil {
		return err
	} else if numUsers == 0 {
		user := NewUser("demo", "demo@photoshare.com", "demo1")
		if err := user.Save(); err != nil {
			return err
		}
	}

	return nil
}

func Init() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "photos.db")
	if err != nil {
		return nil, err
	}
	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	if err := populateDatabase(); err != nil {
		return nil, err
	}

	return db, nil
}
