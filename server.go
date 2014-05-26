package main

import (
	"net/http"
	"os"
)

func main() {

	// STATIC FILES

	http.Handle("/", http.FileServer(http.Dir("./app/")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.ListenAndServe(":"+port, nil)
}
