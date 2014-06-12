package session

import (
	"errors"
	"github.com/danjac/photoshare/api/models"
	"net/http"
	"strconv"
)

const UserCookieName = "userid"

var (
	MissingLoginFields = errors.New("Missing login fields")
	userMgr            = models.NewUserManager()
)

type Authenticator struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func (auth *Authenticator) Identify() (*models.User, error) {

	if auth.Identifier == "" || auth.Password == "" {
		return nil, MissingLoginFields
	}
	return userMgr.Authenticate(auth.Identifier, auth.Password)
}

func GetCurrentUser(r *http.Request) (*models.User, error) {

	userID, err := cookieReader.Read(r, UserCookieName, true)
	if err != nil {
		return nil, err
	}
	if userID == "" {
		return nil, nil
	}
	return userMgr.GetActive(userID)
}

func Login(w http.ResponseWriter, user *models.User) error {
	return cookieWriter.Write(w, UserCookieName, strconv.FormatInt(user.ID, 10), true)
}

func Logout(w http.ResponseWriter) error {
	return cookieWriter.Write(w, UserCookieName, "", true)
}
