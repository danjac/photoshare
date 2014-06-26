package models

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var dbMap *gorp.DbMap

func InitDB(db *sql.DB) (*gorp.DbMap, error) {
	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbMap.TraceOn("[sql]", log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds))

	dbMap.AddTableWithName(User{}, "users").SetKeys(true, "ID")
	dbMap.AddTableWithName(Photo{}, "photos").SetKeys(true, "ID")
	dbMap.AddTableWithName(Tag{}, "tags").SetKeys(true, "ID")

	return dbMap, nil
}
