package render

import (
	"encoding/json"
	"github.com/gorilla/feeds"
	"log"
	"net/http"
	"strconv"
)

func writeBody(w http.ResponseWriter, body []byte, status int, contentType string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", contentType+"; charset=utf8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Write(body)
}

func Atom(w http.ResponseWriter, feed *feeds.Feed) {
	atom, err := feed.ToAtom()
	if err != nil {
		ServerError(w, err)
		return
	}
	writeBody(w, []byte(atom), http.StatusOK, "application/atom+xml")
}

// write a plain text message
func String(w http.ResponseWriter, body string, status int) {
	writeBody(w, []byte(body), status, "text/plain")
}

func ServerError(w http.ResponseWriter, err error) {
	// maybe send email etc in production...
	log.Println(err)
	Error(w, http.StatusInternalServerError)
}

func JSON(w http.ResponseWriter, value interface{}, status int) {
	body, err := json.Marshal(value)
	if err != nil {
		ServerError(w, err)
		return
	}
	writeBody(w, body, status, "application/json")
}

func Error(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func Status(w http.ResponseWriter, status int) {
	String(w, http.StatusText(status), status)
}
