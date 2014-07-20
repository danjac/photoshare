package api

import (
	"errors"
	"github.com/danryan/env"
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
	SmtpPort          int    `env:"key=SMTP_PORT default=25"`
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

func NewAppConfig() (*AppConfig, error) {

	config := &AppConfig{}

	if err := env.Process(config); err != nil {
		return nil, err
	}

	if config.TestDBName == "" {
		config.TestDBName = config.DBName + "_test"
	}

	if config.TestDBUser == "" {
		config.TestDBUser = config.DBUser
	}

	if config.TestDBPassword == "" {
		config.TestDBPassword = config.DBPassword
	}

	if config.TestDBHost == "" {
		config.TestDBHost = config.DBHost
	}

	if config.TestDBName == config.DBName {
		errors.New("Test DB name same as DB name")
	}

	if config.BaseDir == "" {
		config.BaseDir = getDefaultBaseDir()
	}

	if config.PublicDir == "" {
		config.PublicDir = path.Join(config.BaseDir, "public")
	}

	if config.UploadsDir == "" {
		config.UploadsDir = path.Join(config.PublicDir, "uploads")
	}

	if config.ThumbnailsDir == "" {
		config.ThumbnailsDir = path.Join(config.UploadsDir, "thumbnails")
	}

	if config.TemplatesDir == "" {
		config.TemplatesDir = path.Join(config.BaseDir, "templates")
	}

	return config, nil
}

func getDefaultBaseDir() string {
	defaultBaseDir, err := os.Getwd()
	if err != nil {
		defaultBaseDir = "."
	}
	return defaultBaseDir
}
