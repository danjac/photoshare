package api

import (
	"github.com/danryan/env"
	"log"
	"os"
	"path"
)

type AppConfig struct {
	DBName     string `env:"key=DB_NAME required=true"`
	DBUser     string `env:"key=DB_USER required=true"`
	DBPassword string `env:"key=DB_PASS required=true"`
	DBHost     string `env:"key=DB_HOST default=localhost"`

	TestDBName     string `env:"key=TEST_DB_NAME"`
	TestDBUser     string `env:"key=TEST_DB_USER"`
	TestDBPassword string `env:"key=TEST_DB_PASS"`
	TestDBHost     string `env:"key=TEST_DB_HOST"`

	LogSql bool `env:"key=LOG_SQL default=false"`

	SmtpName          string `env:"key=SMTP_NAME"`
	SmtpPassword      string `env:"key=SMTP_PASS"`
	SmtpUser          string `env:"key=SMTP_USER"`
	SmtpHost          string `env:"key=SMTP_HOST default=localhost"`
	SmtpDefaultSender string `env:"key=DEFAULT_EMAIL_SENDER default=webmaster@localhost"`

	BaseDir       string `env:"key=BASE_DIR"`
	PublicDir     string `env:"key=PUBLIC_DIR"`
	UploadsDir    string `env:"key=UPLOADS_DIR"`
	ThumbnailsDir string `env:"key=THUMBNAILS_DIR"`
	TemplatesDir  string `env:"key=TEMPLATES_DIR"`

	PrivateKey string `env:"key=PRIVATE_KEY required=true"`
	PublicKey  string `env:"key=PUBLIC_KEY required=true"`

	ServerPort int `env:"key=PORT default=5000"`
}

var Config = &AppConfig{}

func initConfig() {

	if err := env.Process(Config); err != nil {
		log.Fatal(err)
	}

	if Config.TestDBName == "" {
		Config.TestDBName = Config.DBName + "_test"
	}

	if Config.TestDBUser == "" {
		Config.TestDBUser = Config.DBUser
	}

	if Config.TestDBPassword == "" {
		Config.TestDBPassword = Config.DBPassword
	}

	if Config.TestDBHost == "" {
		Config.TestDBHost = Config.DBHost
	}

	if Config.TestDBName == Config.DBName {
		log.Fatal("Test DB name same as DB name")
	}

	if Config.BaseDir == "" {
		Config.BaseDir = getDefaultBaseDir()
	}

	if Config.PublicDir == "" {
		Config.PublicDir = path.Join(Config.BaseDir, "public")
	}

	if Config.UploadsDir == "" {
		Config.UploadsDir = path.Join(Config.PublicDir, "uploads")
	}

	if Config.ThumbnailsDir == "" {
		Config.ThumbnailsDir = path.Join(Config.UploadsDir, "thumbnails")
	}

	if Config.TemplatesDir == "" {
		Config.TemplatesDir = path.Join(Config.BaseDir, "templates")
	}
}

func getDefaultBaseDir() string {
	defaultBaseDir, err := os.Getwd()
	if err != nil {
		defaultBaseDir = "."
	}
	return defaultBaseDir
}
