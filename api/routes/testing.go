package routes

import (
	"encoding/json"
	"github.com/zenazn/goji/web"
	"net/http/httptest"
)

// unit test helper functions

func parseJsonBody(res *httptest.ResponseRecorder, value interface{}) error {
	return json.Unmarshal([]byte(res.Body.String()), value)
}

func newContext() web.C {
	c := web.C{}
	c.Env = make(map[string]interface{})
	c.URLParams = make(map[string]string)
	return c
}
