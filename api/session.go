package api

import (
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

type SessionManager interface {
	GetCurrentUser(*http.Request) (*User, error)
	Login(http.ResponseWriter, *User) (string, error)
	Logout(http.ResponseWriter) (string, error)
}

func NewSessionManager(config *AppConfig, userMgr UserManager) (SessionManager, error) {
	mgr := &defaultSessionManager{
		config:  config,
		userMgr: userMgr,
	}
	var err error
	mgr.signKey, err = ioutil.ReadFile(config.PrivateKey)
	if err != nil {
		return mgr, err
	}
	mgr.verifyKey, err = ioutil.ReadFile(config.PublicKey)
	if err != nil {
		return mgr, err
	}
	return mgr, nil
}

type defaultSessionManager struct {
	config             *AppConfig
	userMgr            UserManager
	verifyKey, signKey []byte
}

func (mgr *defaultSessionManager) GetCurrentUser(r *http.Request) (*User, error) {

	userID, err := readToken(mgr.verifyKey, r)
	if err != nil {
		return nil, err
	}

	// no token found, user not yet auth'd. Return unauthenticated user

	if userID == 0 {
		return &User{}, nil
	}

	user, err := mgr.userMgr.GetActive(userID)
	if err != nil {
		return user, err
	}
	user.IsAuthenticated = true

	return user, nil
}

func (mgr *defaultSessionManager) Login(w http.ResponseWriter, user *User) (string, error) {
	return writeToken(mgr.signKey, w, user.ID)
}

func (mgr *defaultSessionManager) Logout(w http.ResponseWriter) (string, error) {
	return writeToken(mgr.signKey, w, 0)
}

func readToken(verifyKey []byte, r *http.Request) (int64, error) {
	tokenString := r.Header.Get(tokenHeader)
	if tokenString == "" {
		return 0, nil
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) ([]byte, error) {
		return verifyKey, nil
	})
	switch err.(type) {
	case nil:
		if !token.Valid {
			return 0, nil
		}
		token := token.Claims["uid"].(string)
		if userID, err := strconv.ParseInt(token, 10, 0); err != nil {
			return 0, nil
		} else {
			return userID, nil
		}
	case *jwt.ValidationError:
		return 0, nil
	default:
		return 0, err
	}
}

func writeToken(signKey []byte, w http.ResponseWriter, userID int64) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["uid"] = strconv.FormatInt(userID, 10)
	token.Claims["exp"] = time.Now().Add(time.Minute * expiry).Unix()
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	w.Header().Set(tokenHeader, tokenString)
	return tokenString, nil
}
