package routes

import (
	"encoding/json"
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/session"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"runtime/debug"
)

type AppContext struct {
	*http.Request
	Response http.ResponseWriter
	Params   map[string]string
	User     *models.User
	PhotoMgr models.IPhotoManager
}

func (c *AppContext) Param(name string) string {
	return c.Params[name]
}

func (c *AppContext) GetCurrentUser() (*models.User, error) {
	var err error
	c.User, err = session.GetCurrentUser(c.Request)
	return c.User, err
}

func (c *AppContext) Login(user *models.User) error {
	c.User = user
	return session.Login(c.Response, user)
}

func (c *AppContext) Logout() error {
	c.User = nil
	return session.Logout(c.Response)
}

func (c *AppContext) Render(status int, value interface{}) error {
	c.Response.WriteHeader(status)
	c.Response.Header().Set("Content-type", "application/json")
	return json.NewEncoder(c.Response).Encode(value)
}

func (c *AppContext) OK(value interface{}) error {
	return c.Render(http.StatusOK, value)
}

func (c *AppContext) Unauthorized(value interface{}) error {
	return c.Render(http.StatusUnauthorized, value)
}

func (c *AppContext) Forbidden(value interface{}) error {
	return c.Render(http.StatusForbidden, value)
}

func (c *AppContext) BadRequest(value interface{}) error {
	return c.Render(http.StatusBadRequest, value)
}

func (c *AppContext) NotFound(value interface{}) error {
	return c.Render(http.StatusNotFound, value)
}

func (c *AppContext) ParseJSON(value interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(value)
}

func (c *AppContext) Error(err error) {
	stack := debug.Stack()
	log.Printf("ERROR: %s\n%s\n%s", c.Request, err, stack)
	http.Error(c.Response, "Sorry, an error has occurred", http.StatusInternalServerError)
}

func NewAppContext(w http.ResponseWriter, r *http.Request) *AppContext {
	return &AppContext{r, w, mux.Vars(r), nil, models.PhotoManager}
}

func MakeAppHandler(fn AppHandlerFunc, loginRequired bool) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		c := NewAppContext(w, r)
		if loginRequired {
			if user, err := c.GetCurrentUser(); err != nil || user == nil {
				if err != nil {
					c.Error(err)
					return
				}
				c.Unauthorized("You must be logged in")
				return
			}
		}

		if err := fn(c); err != nil {
			c.Error(err)
		}
	}

}

type AppHandlerFunc func(c *AppContext) error
