package session

import (
	"errors"
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
	MissingLoginFields = errors.New("Missing login fields")
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

// Basic user session info
type SessionInfo struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"isAdmin"`
	LoggedIn bool   `json:"loggedIn"`
}

func NewSessionInfo(user *models.User) *SessionInfo {
	if user == nil || user.ID == 0 || !user.IsAuthenticated {
		return &SessionInfo{}
	}

	return &SessionInfo{user.ID, user.Name, user.IsAdmin, true}
}

// Handles user authentication
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

	userID, err := readToken(r)
	if err != nil {
		return nil, err
	}
	if userID == "" {
		return &models.User{}, nil
	}
	return userMgr.GetActive(userID)
}

func Login(w http.ResponseWriter, user *models.User) (string, error) {
	return createToken(w, strconv.FormatInt(user.ID, 10))
}

func Logout(w http.ResponseWriter) (string, error) {
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
