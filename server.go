package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/routes"
	"github.com/danjac/photoshare/api/settings"
	"github.com/zenazn/goji"
	"log"
	"net/http"
	"os"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s",
		settings.DBUser,
		settings.DBName,
		settings.DBPassword,
		settings.DBHost,
	))

	if err != nil {
		panic(err)
	}

	defer db.Close()

	if _, err := models.InitDB(db); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	flag.Set("bind", "localhost:"+port)

	routes.Setup()

	// for local development
	goji.Get("/*", http.FileServer(http.Dir(settings.PublicDir)))
	goji.Serve()

}
