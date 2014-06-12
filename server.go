package main

import (
	"database/sql"
	"fmt"
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/routes"
	"github.com/danjac/photoshare/api/settings"
	"log"
	"net/http"
	"os"
)

func main() {

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s",
		settings.DBUser,
		settings.DBName,
		settings.DBPassword,
	))

	if err != nil {
		panic(err)
	}

	defer db.Close()

	if _, err := models.InitDB(db); err != nil {
		panic(err)
	}

	http.Handle("/", routes.GetHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Println("starting server on port", port)

	http.ListenAndServe(":"+port, nil)
}
