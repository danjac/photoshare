package photoshare

import (
	"errors"
	"github.com/danryan/env"
	"os"
	"path"
)

type config struct {
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

	ServerPort int `env:"key=API_PORT default=5000"`
}

func newConfig() (*config, error) {
	cfg := &config{}
	if err := env.Process(cfg); err != nil {
		return cfg, err
	}

	if cfg.TestDBName == "" {
		cfg.TestDBName = cfg.DBName + "_test"
	}

	if cfg.TestDBUser == "" {
		cfg.TestDBUser = cfg.DBUser
	}

	if cfg.TestDBPassword == "" {
		cfg.TestDBPassword = cfg.DBPassword
	}

	if cfg.TestDBHost == "" {
		cfg.TestDBHost = cfg.DBHost
	}

	if cfg.TestDBName == cfg.DBName {
		return cfg, errors.New("test DB name same as DB name")
	}

	if cfg.BaseDir == "" {
		cfg.BaseDir = getDefaultBaseDir()
	}

	if cfg.PublicDir == "" {
		cfg.PublicDir = path.Join(cfg.BaseDir, "public")
	}

	if cfg.UploadsDir == "" {
		cfg.UploadsDir = path.Join(cfg.PublicDir, "uploads")
	}

	if cfg.ThumbnailsDir == "" {
		cfg.ThumbnailsDir = path.Join(cfg.UploadsDir, "thumbnails")
	}

	if cfg.TemplatesDir == "" {
		cfg.TemplatesDir = path.Join(cfg.BaseDir, "templates")
	}

	return cfg, nil
}

func getDefaultBaseDir() string {
	defaultBaseDir, err := os.Getwd()
	if err != nil {
		defaultBaseDir = "."
	}
	return defaultBaseDir
}
