package session

import (
	"github.com/danjac/photoshare/api/models"
	"net/http"
)

const UserCookieName = "userid"

var userMgr = models.NewUserManager()

func GetCurrentUser(r *http.Request) (*models.User, error) {
	cookie, err := r.Cookie(UserCookieName)
	if err != nil {
		return nil, nil
	}

	var userID int
	if err := sCookie.Decode(UserCookieName, cookie.Value, &userID); err != nil {
		return nil, nil
	}

	if userID == 0 {
		return nil, nil
	}

	return userMgr.GetActive(userID)
}

func Login(w http.ResponseWriter, user *models.User) error {
	return writeSessionCookie(w, user.ID)
}

func Logout(w http.ResponseWriter) error {
	return writeSessionCookie(w, 0)
}

func writeSessionCookie(w http.ResponseWriter, id int64) error {
	// write the user ID to the secure cookie
	encoded, err := sCookie.Encode(UserCookieName, id)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:  UserCookieName,
		Value: encoded,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	return nil

}
