package routes

import (
	"encoding/json"
	"fmt"
	"github.com/danjac/photoshare/api/config"
	"net/http"
	"path"
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

func parseTemplate(name string) *template.Template {
	return template.Must(template.ParseFiles(path.Join(config.Dirs.Templates, name)))
}
