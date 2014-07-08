package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

func init() {
	initConfig()
	initEmail()
	initSession()
}

func getDbConn() (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s",
		config.DBUser,
		config.DBName,
		config.DBPassword,
		config.DBHost,
	))

}

func RunServer() {

	runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)

	db, err := getDbConn()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if _, err := InitDB(db, config.LogSql); err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(fmt.Sprintf("localhost:%d", config.ServerPort), setupRoutes())

}
