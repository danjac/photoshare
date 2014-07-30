package photoshare

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/coopernurse/gorp"
	"net/http/httptest"
)

// unit test helper functions

func parseJSONBody(res *httptest.ResponseRecorder, value interface{}) error {
	return json.Unmarshal([]byte(res.Body.String()), value)
}

type testDB struct {
	DB    *sql.DB
	dbMap *gorp.DbMap
}

func (tdb *testDB) clean() {
	var tables = []string{"photo_tags", "tags", "photos", "users"}
	for _, table := range tables {
		if _, err := tdb.dbMap.Exec("DELETE FROM " + table); err != nil {
			panic(err)
		}
	}
	defer tdb.DB.Close()
}

func makeTestDB(cfg *appConfig) (tdb *testDB) {
	var err error

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s",
		cfg.TestDBUser,
		cfg.TestDBName,
		cfg.TestDBPassword,
		cfg.TestDBHost,
	))

	if err != nil {
		panic(err)
	}

	dbMap, err := initDB(db, false)
	if err != nil {
		panic(err)
	}

	return &testDB{db, dbMap}
}
