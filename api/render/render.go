package render

import (
	"encoding/json"
	"github.com/gorilla/feeds"
	"net/http"
	"strconv"
)

func writeBody(w http.ResponseWriter, body []byte, status int, contentType string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", contentType+";charset=utf8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Write(body)
}

// write a plain text message
func String(w http.ResponseWriter, body string, status int) {
	writeBody(w, []byte(body), status, "text/plain")
}

func JSON(w http.ResponseWriter, value interface{}, status int) {
	body, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	writeBody(w, body, status, "application/json")
}

func Atom(w http.ResponseWriter, feed *feeds.Feed, status int) {
	atom, err := feed.ToAtom()
	if err != nil {
		panic(err)
	}
	writeBody(w, []byte(atom), status, "application/atom+xml")
}

func Error(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
