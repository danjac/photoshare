package models

import (
	"database/sql"
	"fmt"
	"github.com/danjac/photoshare/api/settings"
)

type TestDB struct {
	DB *sql.DB
}

func (tdb *TestDB) Clean() {
    /*
	if err := dbMap.TruncateTables(); err != nil {
		panic(err)
	}
    */
    var tables = []string{"photo_tags", "tags", "photos", "users"}
    for _, table := range tables {
        if _, err := dbMap.Exec("DELETE FROM " + table); err != nil {
            panic(err)
        }
    }
	defer tdb.DB.Close()
}

func MakeTestDB() (tdb *TestDB) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s",
		settings.TestDBUser,
		settings.TestDBName,
		settings.TestDBPassword,
	))

	if err != nil {
		panic(err)
	}
	if _, err := InitDB(db); err != nil {
		panic(err)
	}

	return &TestDB{db}
}
