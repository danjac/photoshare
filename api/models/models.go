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
	dbMap.TraceOn("[sql]", log.New(os.Stdout, "photoshare:", log.Lmicroseconds))
	return dbMap, nil
}
