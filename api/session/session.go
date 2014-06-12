package session

import (
	"github.com/danjac/photoshare/api/settings"
	"github.com/gorilla/securecookie"
	"net/http"
	"time"
)

type CookieReader interface {
	Read(*http.Request, string, bool) (string, error)
}

type CookieWriter interface {
	Write(http.ResponseWriter, string, string, bool) error
}

type DefaultCookieReader struct {
	*securecookie.SecureCookie
}

func (reader *DefaultCookieReader) Read(r *http.Request, name string, decode bool) (string, error) {

	var value string

	cookie, err := r.Cookie(name)

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
		if err := reader.Decode(name, cookie.Value, &value); err != nil {
			return value, err
		}
	} else {
		value = cookie.Value
	}

	return value, nil
}

type DefaultCookieWriter struct {
	*securecookie.SecureCookie
}

func (writer *DefaultCookieWriter) Write(w http.ResponseWriter, name string, value string, encode bool) error {

	if encode {
		var err error
		value, err = writer.Encode(name, value)
		if err != nil {
			return err
		}
	}

	expires := time.Now().AddDate(0, 0, 1)

	cookie := &http.Cookie{
		Name:       name,
		Value:      value,
		Path:       "/",
		MaxAge:     86400,
		Expires:    expires,
		RawExpires: expires.Format(time.UnixDate),
	}
	http.SetCookie(w, cookie)
	return nil
}

var sCookie = securecookie.New([]byte(settings.HashKey),
	[]byte(settings.BlockKey))

var cookieReader = &DefaultCookieReader{sCookie}
var cookieWriter = &DefaultCookieWriter{sCookie}
