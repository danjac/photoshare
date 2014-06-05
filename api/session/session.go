package session

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/gorilla/securecookie"
	"net/http"
)

const (
	CookieName = "userid"
)

var hashKey = securecookie.GenerateRandomKey(32)
var blockKey = securecookie.GenerateRandomKey(32)
var sCookie = securecookie.New(hashKey, blockKey)

func GetCurrentUser(r *http.Request) (*models.User, error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return nil, nil
	}

	var userID int
	if err := sCookie.Decode(CookieName, cookie.Value, &userID); err != nil {
		return nil, nil
	}

	if userID == 0 {
		return nil, nil
	}

	return models.GetUser(userID)
}

func Login(w http.ResponseWriter, user *models.User) error {
	return writeSessionCookie(w, user.ID)
}

func Logout(w http.ResponseWriter) error {
	return writeSessionCookie(w, 0)
}

func writeSessionCookie(w http.ResponseWriter, id int) error {
	// write the user ID to the secure cookie
	encoded, err := sCookie.Encode(CookieName, id)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:  CookieName,
		Value: encoded,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	return nil

}
