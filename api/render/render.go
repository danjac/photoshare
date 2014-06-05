package render

import (
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func Status(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

func JSON(w http.ResponseWriter, status int, value interface{}) {
	w.WriteHeader(status)
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(value)
}
