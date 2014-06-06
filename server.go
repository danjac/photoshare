package main

import (
	"fmt"
	"github.com/danjac/photoshare/api"
	"github.com/danjac/photoshare/api/settings"
	"log"
	"net/http"
	"os"
)

func getEnvOrDie(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatal(fmt.Sprintf("environ key %s is missing", name))
	}
	return value
}

func getEnvOrElse(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {

	config := &settings.AppConfig{
		DBName:           getEnvOrDie("DB_NAME"),
		DBUser:           getEnvOrDie("DB_USER"),
		DBPassword:       getEnvOrDie("DB_PASS"),
		DBHost:           getEnvOrElse("DB_HOST", "localhost"),
		ApiPathPrefix:    getEnvOrElse("API_PATH", "/api"),
		PublicPathPrefix: getEnvOrElse("PUBLIC_PATH", "/"),
		PublicDir:        getEnvOrElse("PUBLIC_DIR", "./public/"),
		UploadsDir:       getEnvOrElse("UPLOADS_DIR", "public/uploads"),
		LogPrefix:        getEnvOrElse("LOG_PREFIX", "photoshare"),
	}

	app, err := api.NewApplication(config)
	if err != nil {
		log.Fatal(err)
	}
	defer app.Shutdown()

	http.Handle("/", app.Handler)

	port := getEnvOrElse("PORT", "5000")

	log.Println("starting server on port", port)

	http.ListenAndServe(":"+port, nil)
}
