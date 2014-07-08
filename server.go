package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/danjac/photoshare/api"
	"github.com/zenazn/goji"
	"log"
	"net/http"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s",
		api.Config.DBUser,
		api.Config.DBName,
		api.Config.DBPassword,
		api.Config.DBHost,
	))

	if err != nil {
		panic(err)
	}

	defer db.Close()

	if _, err := api.InitDB(db, api.Config.LogSql); err != nil {
		log.Fatal(err)
	}

	flag.Set("bind", fmt.Sprintf("localhost:%d", api.Config.ServerPort))

	// for local development
	goji.Get("/*", http.FileServer(http.Dir(api.Config.PublicDir)))
	goji.Serve()

}
