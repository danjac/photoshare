package settings

import (
	"log"
	"os"
	"path"
)

var (
	PrivKeyFile,
	PubKeyFile,
	DBHost,
	DBName,
	DBUser,
	DBPassword,
	TestDBName,
	TestDBUser,
	TestDBPassword,
	TestDBHost,
	LogPrefix,
	PublicDir,
	UploadsDir,
	ThumbnailsDir string
)

func getEnvOrDie(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatal("Environment setting ", name, " is missing")
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

func init() {

	PrivKeyFile = getEnvOrDie("PRIVATE_KEY")
	PubKeyFile = getEnvOrDie("PUBLIC_KEY")

	LogPrefix = getEnvOrElse("LOG_PREFIX", "photoshare")

	DBName = getEnvOrDie("DB_NAME")
	DBUser = getEnvOrDie("DB_USER")
	DBPassword = getEnvOrDie("DB_PASS")
	DBHost = getEnvOrElse("DB_HOST", "localhost")

	TestDBName = getEnvOrDie("TEST_DB_NAME")
	TestDBUser = getEnvOrDie("TEST_DB_USER")
	TestDBPassword = getEnvOrDie("TEST_DB_PASS")
	TestDBHost = getEnvOrElse("TEST_DB_HOST", "localhost")

	if TestDBName == DBName {
		log.Fatal("Test DB name same as DB name")
	}

	PublicDir = getEnvOrElse("PUBLIC_DIR", "./public/")
	UploadsDir = getEnvOrElse("UPLOADS_DIR", path.Join(PublicDir, "uploads"))
	ThumbnailsDir = getEnvOrElse("THUMBNAILS_DIR", path.Join(UploadsDir, "thumbnails"))
}
