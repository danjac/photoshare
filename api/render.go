package api

import (
	"encoding/json"
	"github.com/gorilla/feeds"
	"log"
	"net/http"
	"strconv"
)

type Render struct{}

var render = &Render{}

func (render *Render) writeBody(w http.ResponseWriter, body []byte, status int, contentType string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", contentType+"; charset=utf8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Write(body)
}

func (render *Render) Atom(w http.ResponseWriter, feed *feeds.Feed) {
	atom, err := feed.ToAtom()
	if err != nil {
		render.ServerError(w, err)
		return
	}
	render.writeBody(w, []byte(atom), http.StatusOK, "application/atom+xml")
}

// write a plain text message
func (render *Render) String(w http.ResponseWriter, body string, status int) {
	render.writeBody(w, []byte(body), status, "text/plain")
}

func (render *Render) ServerError(w http.ResponseWriter, err error) {
	// maybe send email etc in production...
	log.Println(err)
	render.Error(w, http.StatusInternalServerError)
}

func (render *Render) JSON(w http.ResponseWriter, value interface{}, status int) {
	body, err := json.Marshal(value)
	if err != nil {
		render.ServerError(w, err)
		return
	}
	render.writeBody(w, body, status, "application/json")
}

func (render *Render) Error(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (render *Render) Status(w http.ResponseWriter, status int) {
	render.String(w, http.StatusText(status), status)
}
