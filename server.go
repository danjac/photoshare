package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/danjac/photoshare/api/config"
	"github.com/danjac/photoshare/api/models"
	_ "github.com/danjac/photoshare/api/routes"
	"github.com/zenazn/goji"
	"log"
	"net/http"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s",
		config.DBUser,
		config.DBName,
		config.DBPassword,
		config.DBHost,
	))

	if err != nil {
		panic(err)
	}

	defer db.Close()

	if _, err := models.InitDB(db, config.LogSql); err != nil {
		log.Fatal(err)
	}

	flag.Set("bind", "localhost:"+config.ServerPort)

	// for local development
	goji.Get("/*", http.FileServer(http.Dir(config.PublicDir)))
	goji.Serve()

}
