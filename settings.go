package photoshare

import (
	"errors"
	"github.com/danryan/env"
	"os"
	"path"
)

type settings struct {
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

	MemcacheHost string `env:"key=MEMCACHE_HOST default=0.0.0.0:11211"`

	GoogleClientID string `env:"key=GOOGLE_CLIENT_ID"`
	GoogleSecret   string `env:"key=GOOGLE_SECRET"`

	ServerPort int `env:"key=PORT default=5000"`
}

func newSettings() (*settings, error) {
	s := &settings{}

	if err := env.Process(s); err != nil {
		return s, err
	}

	if s.TestDBName == "" {
		s.TestDBName = s.DBName + "_test"
	}

	if s.TestDBUser == "" {
		s.TestDBUser = s.DBUser
	}

	if s.TestDBPassword == "" {
		s.TestDBPassword = s.DBPassword
	}

	if s.TestDBHost == "" {
		s.TestDBHost = s.DBHost
	}

	if s.TestDBName == s.DBName {
		return s, errors.New("test DB name same as DB name")
	}

	if s.BaseDir == "" {
		s.BaseDir = getDefaultBaseDir()
	}

	if s.PublicDir == "" {
		s.PublicDir = path.Join(s.BaseDir, "public")
	}

	if s.UploadsDir == "" {
		s.UploadsDir = path.Join(s.PublicDir, "uploads")
	}

	if s.ThumbnailsDir == "" {
		s.ThumbnailsDir = path.Join(s.UploadsDir, "thumbnails")
	}

	if s.TemplatesDir == "" {
		s.TemplatesDir = path.Join(s.BaseDir, "templates")
	}

	return s, nil
}

func getDefaultBaseDir() string {
	defaultBaseDir, err := os.Getwd()
	if err != nil {
		defaultBaseDir = "."
	}
	return defaultBaseDir
}
