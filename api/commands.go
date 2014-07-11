package api

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/zenazn/goji"
	"log"
	"runtime"
)

func init() {
	initConfig()
	initEmail()
	initSession()
	initRoutes()
}

func getDbConn() (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s",
		config.DBUser,
		config.DBName,
		config.DBPassword,
		config.DBHost,
	))

}

func Serve() {

	runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)

	db, err := getDbConn()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if _, err := InitDB(db, config.LogSql); err != nil {
		log.Fatal(err)
	}

	flag.Set("bind", fmt.Sprintf(":%d", config.ServerPort))

	goji.Serve()

}
