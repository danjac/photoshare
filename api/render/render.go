package render

import (
	"encoding/json"
	"log"
	"net/http"
)

func Error(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err, r)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func Status(w http.ResponseWriter, status int, msg string) error {
	w.WriteHeader(status)
	w.Write([]byte(msg))
	return nil
}

func JSON(w http.ResponseWriter, status int, value interface{}) error {
	w.WriteHeader(status)
	w.Header().Add("content-type", "application/json")
	return json.NewEncoder(w).Encode(value)
}
