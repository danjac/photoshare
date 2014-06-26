package session

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/settings"
	jwt "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	tokenHeader = "X-Auth-Token"
	expiry      = 60 // minutes
)

var (
	verifyKey, signKey []byte
	userMgr            = models.NewUserManager()
)

func init() {
	var err error
	signKey, err = ioutil.ReadFile(settings.PrivKeyFile)
	if err != nil {
		panic(err)
	}
	verifyKey, err = ioutil.ReadFile(settings.PubKeyFile)
	if err != nil {
		panic(err)
	}

}

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

func readToken(r *http.Request) (string, error) {
	tokenString := r.Header.Get(tokenHeader)
	if tokenString == "" {
		return "", nil
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) ([]byte, error) {
		return verifyKey, nil
	})
	switch err.(type) {
	case nil:
		if !token.Valid {
			return "", nil
		}
		return token.Claims["uid"].(string), nil
	case *jwt.ValidationError:
		return "", nil
	default:
		return "", err
	}
}

func createToken(w http.ResponseWriter, userID string) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["uid"] = userID
	token.Claims["exp"] = time.Now().Add(time.Minute * expiry).Unix()
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	w.Header().Set(tokenHeader, tokenString)
	return tokenString, nil
}
