package models

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/danjac/photoshare/api/settings"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var dbMap *gorp.DbMap

func populateDatabase() error {

	log.Println("Populating database....")

	numUsers, err := dbMap.SelectInt("SELECT COUNT(id) FROM users")
	if err != nil {
		return err
	} else if numUsers == 0 {
		log.Println("CREATING DEMO USER")
		user := &User{Name: "demo", Email: "demo@photoshare.com", Password: "demo1"}
		if err := user.Insert(); err != nil {
			return err
		}
	}

	return nil
}

func Init() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s",
		settings.Config.DBUser,
		settings.Config.DBName,
		settings.Config.DBPassword))
	if err != nil {
		return nil, err
	}

	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbMap.TraceOn("[sql]", log.New(os.Stdout, settings.Config.LogPrefix+":", log.Lmicroseconds))

	dbMap.AddTableWithName(User{}, "users").SetKeys(true, "ID")
	dbMap.AddTableWithName(Photo{}, "photos").SetKeys(true, "ID")
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		return nil, err
	}

	if err := populateDatabase(); err != nil {
		return nil, err
	}

	return db, nil
}
