package models

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
)

var dbMap *gorp.DbMap

func InitDB(db *sql.DB) (*gorp.DbMap, error) {
	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	//dbMap.TraceOn("[sql]", log.New(os.Stdout, "photoshare:", log.Lmicroseconds))

	dbMap.AddTableWithName(User{}, "users").SetKeys(true, "ID")
	dbMap.AddTableWithName(Photo{}, "photos").SetKeys(true, "ID")
	dbMap.AddTableWithName(Tag{}, "tags").SetKeys(true, "ID")

	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		return dbMap, err
	}
	return dbMap, nil
}
