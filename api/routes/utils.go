package routes

import (
	"encoding/json"
	"fmt"
	"github.com/danjac/photoshare/api/config"
	"net/http"
	"path"
	"strconv"
	"text/template"
)

func parseJSON(r *http.Request, value interface{}) error {
	return json.NewDecoder(r.Body).Decode(value)
}

func scheme(r *http.Request) string {
	if r.TLS == nil {
		return "http"
	}
	return "https"
}

func baseURL(r *http.Request) string {
	return fmt.Sprintf("%s://%s", scheme(r), r.Host)
}

func writeBody(w http.ResponseWriter, body []byte, status int, contentType string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", contentType+";charset=utf8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Write(body)
}

// write a plain text message
func writeString(w http.ResponseWriter, body string, status int) {
	writeBody(w, []byte(body), status, "text/plain")
}

func writeJSON(w http.ResponseWriter, value interface{}, status int) {
	body, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	writeBody(w, body, status, "application/json")
}

func writeError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func parseTemplate(name string) *template.Template {
	return template.Must(template.ParseFiles(path.Join(config.TemplateDir, name)))
}
