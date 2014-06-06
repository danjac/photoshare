package api

import (
	"database/sql"
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/routes"
	"github.com/danjac/photoshare/api/settings"
	"net/http"
)

type Application struct {
	DB      *sql.DB
	Handler http.Handler
}

func (app *Application) Shutdown() {
	app.DB.Close()
}

func NewApplication(config *settings.AppConfig) (*Application, error) {

	settings.Config = config
	app := &Application{}

	db, err := models.Init()
	if err != nil {
		return nil, err
	}
	app.DB = db

	app.Handler = routes.Init()

	return app, nil
}
