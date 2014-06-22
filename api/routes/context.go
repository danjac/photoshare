package routes

import (
	"encoding/json"
	"fmt"
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/session"
	"github.com/danjac/photoshare/api/settings"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type Result struct {
	http.ResponseWriter
	Status      int
	Body        []byte
	ContentType string
	Error       error
}

func (r *Result) Render() error {
	if r.Error != nil {
		http.Error(r.ResponseWriter, string(r.Body), r.Status)
		return r.Error
	}

	r.WriteHeader(r.Status)
	r.Header().Set("Content-type", r.ContentType)
	_, err := r.Write(r.Body)
	return err
}

type Context struct {
	*http.Request
	Response http.ResponseWriter
	Params   map[string]string
	User     *models.User
	Log      *log.Logger
}

func (c *Context) Result(status int, body []byte, contentType string, err error) *Result {
	return &Result{c.Response, status, body, contentType, err}
}

func (c *Context) Json(status int, value interface{}) *Result {
	body, err := json.Marshal(value)
	if err != nil {
		return c.Error(err)
	}
	return c.Result(status, body, "application/json", nil)
}

// Renders feed in Atom format
func (c *Context) Atom(feed *feeds.Feed) *Result {
	atom, err := feed.ToAtom()
	if err != nil {
		return c.Error(err)
	}
	return c.Result(http.StatusOK, []byte(atom), "application/atom+xml", nil)
}

func (c *Context) Param(name string) string {
	return c.Params[name]
}

func (c *Context) Scheme() string {
	if c.Request.TLS == nil {
		return "http"
	}
	return "https"
}

func (c *Context) BaseURL() string {
	return fmt.Sprintf("%s://%s", c.Scheme(), c.Request.Host)
}

func (c *Context) GetCurrentUser() (*models.User, error) {
	var err error
	if c.User != nil {
		return c.User, nil
	}
	c.User, err = session.GetCurrentUser(c.Request)
	return c.User, err
}

func (c *Context) Login(user *models.User) error {
	c.User = user
	_, err := session.Login(c.Response, user)
	return err
}

func (c *Context) Logout() error {
	if c.User != nil {
		c.User.IsAuthenticated = false
	}
	_, err := session.Logout(c.Response)
	return err
}

func (c *Context) OK(value interface{}) *Result {
	return c.Json(http.StatusOK, value)
}

func (c *Context) Unauthorized(value interface{}) *Result {
	return c.Json(http.StatusUnauthorized, value)
}

func (c *Context) Forbidden(value interface{}) *Result {
	return c.Json(http.StatusForbidden, value)
}

func (c *Context) BadRequest(value interface{}) *Result {
	return c.Json(http.StatusBadRequest, value)
}

func (c *Context) NotFound(value interface{}) *Result {
	return c.Json(http.StatusNotFound, value)
}

func (c *Context) Error(err error) *Result {
	return c.Result(http.StatusInternalServerError, []byte("error"), "text/plain", err)
}

func (c *Context) ParseJSON(value interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(value)
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {

	prefix := fmt.Sprintf("%s:[%s] %s ", settings.LogPrefix, r.Method, r.URL.Path)
	logger := log.New(os.Stdout, prefix, log.Lmicroseconds)
	return &Context{r, w, mux.Vars(r), nil, logger}
}
