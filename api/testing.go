package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/zenazn/goji/web"
	"net/http/httptest"
)

// unit test helper functions

func parseJsonBody(res *httptest.ResponseRecorder, value interface{}) error {
	return json.Unmarshal([]byte(res.Body.String()), value)
}

func newContext() web.C {
	c := web.C{}
	c.Env = make(map[string]interface{})
	c.URLParams = make(map[string]string)
	return c
}

type TestDB struct {
	DB    *sql.DB
	dbMap *gorp.DbMap
}

func (tdb *TestDB) Clean() {
	var tables = []string{"photo_tags", "tags", "photos", "users"}
	for _, table := range tables {
		if _, err := tdb.dbMap.Exec("DELETE FROM " + table); err != nil {
			panic(err)
		}
	}
	defer tdb.DB.Close()
}

func MakeTestDB(config *AppConfig) (tdb *TestDB) {
	var err error

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s",
		config.TestDBUser,
		config.TestDBName,
		config.TestDBPassword,
		config.TestDBHost,
	))

	if err != nil {
		panic(err)
	}

	dbMap, err := InitDB(db, false)
	if err != nil {
		panic(err)
	}

	return &TestDB{db, dbMap}
}
