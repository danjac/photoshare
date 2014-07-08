package api

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/zenazn/goji"
	"log"
	"net/http"
	"runtime"
)

func init() {
	initConfig()
	initRoutes()
	initEmail()
	initSession()
}

func getDbConn() (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s",
		Config.DBUser,
		Config.DBName,
		Config.DBPassword,
		Config.DBHost,
	))

}

func RunServer() {

	runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)

	db, err := getDbConn()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if _, err := InitDB(db, Config.LogSql); err != nil {
		log.Fatal(err)
	}

	flag.Set("bind", fmt.Sprintf("localhost:%d", Config.ServerPort))

	// for local development
	goji.Get("/*", http.FileServer(http.Dir(Config.PublicDir)))
	goji.Serve()

}
