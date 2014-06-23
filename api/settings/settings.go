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

	TestDBName = getEnvOrElse("TEST_DB_NAME", DBName+"_test")
	TestDBUser = getEnvOrElse("TEST_DB_USER", DBUser)
	TestDBPassword = getEnvOrElse("TEST_DB_PASS", DBPassword)
	TestDBHost = getEnvOrElse("TEST_DB_HOST", DBHost)

	if TestDBName == DBName {
		log.Fatal("Test DB name same as DB name")
	}

	PublicDir = getEnvOrElse("PUBLIC_DIR", "./public/")
	UploadsDir = getEnvOrElse("UPLOADS_DIR", path.Join(PublicDir, "uploads"))
	ThumbnailsDir = getEnvOrElse("THUMBNAILS_DIR", path.Join(UploadsDir, "thumbnails"))
}
