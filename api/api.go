package api

import (
	"database/sql"
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/routes"
	"net/http"
)

type Config struct {
	DBHost, DBName, DBUser, DBPassword, LogPrefix, UploadsDir, ApiPathPrefix, PublicPathPrefix, PublicDir string
}

type Application struct {
	DB      *sql.DB
	Handler http.Handler
}

func (app *Application) Shutdown() {
	app.DB.Close()
}

func NewApplication(config *Config) (*Application, error) {

	app := &Application{}
	db, err := models.Init(config.DBName,
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.LogPrefix)
	if err != nil {
		return nil, err
	}
	app.DB = db

	app.Handler = routes.Init(config.UploadsDir,
		config.ApiPathPrefix,
		config.PublicPathPrefix,
		config.PublicDir)

	return app, nil
}
