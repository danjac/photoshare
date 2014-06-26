package session

import (
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

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
