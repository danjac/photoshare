package main

import (
	"fmt"
	"github.com/danjac/photoshare/api"
	"net/http"
	"os"
)

func main() {
	app, err := api.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.Shutdown()

	http.Handle("/", app.Handler)

	fmt.Println("starting server...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	http.ListenAndServe(":"+port, nil)
}
