package session

import (
	"github.com/danjac/photoshare/api/settings"
	"github.com/gorilla/securecookie"
	"net/http"
	"time"
)

type CookieManager struct {
	*securecookie.SecureCookie
}

func (mgr *CookieManager) Read(r *http.Request, cookieName string, decode bool) (string, error) {

	var value string

	cookie, err := r.Cookie(cookieName)

	if err != nil {
		if err == http.ErrNoCookie {
			return value, nil
		}
		return value, err
	}

	if cookie.Value == "" {
		return value, nil
	}

	if decode {
		if err := mgr.Decode(cookieName, cookie.Value, &value); err != nil {
			return value, err
		}
	} else {
		value = cookie.Value
	}

	return value, nil
}

func (mgr *CookieManager) Write(w http.ResponseWriter, cookieName string, value string, encode bool) error {

	if encode {
		var err error
		value, err = mgr.Encode(cookieName, value)
		if err != nil {
			return err
		}
	}

	expires := time.Now().AddDate(0, 0, 1)

	cookie := &http.Cookie{
		Name:       cookieName,
		Value:      value,
		Path:       "/",
		MaxAge:     86400,
		Expires:    expires,
		RawExpires: expires.Format(time.UnixDate),
	}
	http.SetCookie(w, cookie)
	return nil
}

func NewCookieManager(hashKey, blockKey string) *CookieManager {
	return &CookieManager{
		securecookie.New([]byte(hashKey),
			[]byte(blockKey)),
	}
}

var cookieMgr = NewCookieManager(settings.HashKey, settings.BlockKey)
