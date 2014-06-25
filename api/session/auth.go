package session

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/settings"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/zenazn/goji/web"
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

// lazily checks for current user in session

func GetCurrentUser(c web.C, r *http.Request) (*models.User, error) {

	obj, ok := c.Env["user"]
	if ok {
		return obj.(*models.User), nil
	}

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

	c.Env["user"] = user
	return user, nil
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
