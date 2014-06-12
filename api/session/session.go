package session

import (
	"github.com/danjac/photoshare/api/settings"
	"github.com/gorilla/securecookie"
)

var sCookie = securecookie.New([]byte(settings.HashKey), []byte(settings.BlockKey))
