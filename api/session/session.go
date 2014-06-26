package session

import (
	"github.com/danjac/photoshare/api/models"
	"net/http"
	"strconv"
)

var userMgr = models.NewUserManager()

type SessionManager interface {
	GetCurrentUser(r *http.Request) (*models.User, error)
	Login(w http.ResponseWriter, user *models.User) (string, error)
	Logout(w http.ResponseWriter) (string, error)
}

func NewSessionManager() SessionManager {
	return &defaultSessionManager{}
}

type defaultSessionManager struct{}

func (mgr *defaultSessionManager) GetCurrentUser(r *http.Request) (*models.User, error) {

	userID, err := readToken(r)
	if err != nil {
		return nil, err
	}

	// no token found, user not yet auth'd. Return unauthenticated user

	if userID == "" {
		return &models.User{}, nil
	}

	user, err := userMgr.GetActive(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (mgr *defaultSessionManager) Login(w http.ResponseWriter, user *models.User) (string, error) {
	return createToken(w, strconv.FormatInt(user.ID, 10))
}

func (mgr *defaultSessionManager) Logout(w http.ResponseWriter) (string, error) {
	return createToken(w, "")
}
