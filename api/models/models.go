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

func InitMap(db *sql.DB, logPrefix string) (*gorp.DbMap, error) {
	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	if logPrefix != "" {
		dbMap.TraceOn("[sql]", log.New(os.Stdout, logPrefix+":", log.Lmicroseconds))
	}

	dbMap.AddTableWithName(User{}, "users").SetKeys(true, "ID")
	dbMap.AddTableWithName(Photo{}, "photos").SetKeys(true, "ID")
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		return dbMap, err
	}
	return dbMap, nil
}

func Init() (*sql.DB, error) {

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s",
		settings.Config.DBUser,
		settings.Config.DBName,
		settings.Config.DBPassword,
	))
	if err != nil {
		return nil, err
	}

	if _, err := InitMap(db, settings.Config.LogPrefix); err != nil {
		return db, err
	}
	return db, nil
}
