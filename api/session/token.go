package session

import (
	"github.com/danjac/photoshare/api/config"
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
)

func init() {
	var err error
	signKey, err = ioutil.ReadFile(config.Keys.Private)
	if err != nil {
		panic(err)
	}
	verifyKey, err = ioutil.ReadFile(config.Keys.Public)
	if err != nil {
		panic(err)
	}

}

func readToken(r *http.Request) (int64, error) {
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

func writeToken(w http.ResponseWriter, userID int64) (string, error) {
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
