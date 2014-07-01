package models

import (
	"database/sql"
	"fmt"
	"github.com/danjac/photoshare/api/config"
)

type TestDB struct {
	DB *sql.DB
}

func (tdb *TestDB) Clean() {
	var tables = []string{"photo_tags", "tags", "photos", "users"}
	for _, table := range tables {
		if _, err := dbMap.Exec("DELETE FROM " + table); err != nil {
			panic(err)
		}
	}
	defer tdb.DB.Close()
}

func MakeTestDB() (tdb *TestDB) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s",
		config.TestDB.User,
		config.TestDB.Name,
		config.TestDB.Password,
		config.TestDB.Host,
	))

	if err != nil {
		panic(err)
	}
	if _, err := InitDB(db, false); err != nil {
		panic(err)
	}

	return &TestDB{db}
}
