package routes

import (
	"log"
	"net/http"
)

type AppHandlerFunc func(c *Context) *Result

func NewAppHandler(fn AppHandlerFunc, loginRequired bool) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var result *Result

		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()

		// set common headers

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
			w.Header().Set("Access-Control-Max-Age", "1000")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Access-Token, X-Requested-With, Content-Type, Accept")
			return
		}

		// new Context

		c := NewContext(w, r)

		if loginRequired {
			if user, err := c.GetCurrentUser(); err != nil || !user.IsAuthenticated {
				if err != nil {
					result = c.Error(err)
				}
				result = c.Unauthorized("You must be logged in")
			}
		}

		if result == nil {
			result = fn(c)
		}

		if err := result.Render(); err != nil {
			c.Log.Panic(err)
		}
	}

}
