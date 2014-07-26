package photoshare

import (
	"database/sql"
	"fmt"
	"github.com/juju/errgo"
	"log"
	"net/http"
)

type httpError struct {
	Status      int
	Description string
}

func (h httpError) Error() string {
	if h.Description == "" {
		return http.StatusText(h.Status)
	}
	return h.Description
}

func isErrSqlNoRows(err error) bool {
	if err == sql.ErrNoRows {
		return true
	}
	if err, ok := err.(*errgo.Err); ok && err.Underlying() == sql.ErrNoRows {
		return true
	}
	return false
}

func logError(err error) {
	s := fmt.Sprintf("Error:%s", err)
	if err, ok := err.(errgo.Locationer); ok {
		s += fmt.Sprintf(" %s", err.Location())
	}
	log.Println(s)
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	if err, ok := err.(httpError); ok {
		http.Error(w, err.Error(), err.Status)
		return
	}

	if err, ok := err.(validationFailure); ok {
		renderJSON(w, err, http.StatusBadRequest)
		return
	}

	if isErrSqlNoRows(err) {
		http.NotFound(w, r)
		return
	}

	logError(err)

	http.Error(w, "Sorry, an error occurred", http.StatusInternalServerError)
}
