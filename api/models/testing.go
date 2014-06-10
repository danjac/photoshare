package models

import (
	"database/sql"
	"fmt"
	"os"
)

type TestDB struct {
	DB *sql.DB
}

func (tdb *TestDB) Clean() {
	if err := dbMap.TruncateTables(); err != nil {
		panic(err)
	}
	defer tdb.DB.Close()
}

func MakeTestDB() (tdb *TestDB) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s",
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_NAME"),
		os.Getenv("TEST_DB_PASS"),
	))

	if err != nil {
		panic(err)
	}
	if _, err := InitMap(db, "TESTING"); err != nil {
		panic(err)
	}

	return &TestDB{db}
}


