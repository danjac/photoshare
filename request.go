package photoshare

import (
	"encoding/json"
	"fmt"
	"github.com/juju/errgo"
	"net/http"
	"strconv"
)

type request struct {
	*http.Request
}

func (r *request) scheme() string {
	if r.TLS == nil {
		return "http"
	}
	return "https"
}

func (r *request) baseURL() string {
	return fmt.Sprintf("%s://%s", r.scheme(), r.Host)
}

func (r *request) decodeJSON(value interface{}) error {
	return errgo.Mask(json.NewDecoder(r.Body).Decode(value))
}

func (r *request) getPage() *page {
	pageNum, err := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if err != nil {
		pageNum = 1
	}
	return newPage(pageNum)
}
