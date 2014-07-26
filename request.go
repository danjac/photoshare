package photoshare

import (
	"encoding/json"
	"fmt"
	"github.com/juju/errgo"
	"github.com/zenazn/goji/web"
	"net/http"
	"strconv"
)

type request struct {
	*http.Request
	c    web.C
	user *user
}

func newRequest(c web.C, r *http.Request) *request {
	req := &request{Request: r}
	req.c = c
	return req
}

func (r *request) getIntParam(name string) int64 {
	value, _ := strconv.ParseInt(r.c.URLParams[name], 10, 0)
	return value
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
