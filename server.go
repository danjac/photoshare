package main

import (
	"fmt"
    "github.com/danjac/photoshare/api/models"
    "github.com/danjac/photoshare/api/routes"
	"net/http"
	"os"
)

func main() {
	db, err := models.Init()

	if err != nil {
		panic(err)
	}
	defer db.Close()

    r := routes.Init()
	http.Handle("/", r)

	fmt.Println("starting server...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	http.ListenAndServe(":"+port, nil)
}
