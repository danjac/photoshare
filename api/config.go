package api

import (
	"github.com/danryan/env"
	"log"
	"os"
	"path"
)

type DBConfig struct {
	Name     string `env:"key=DB_NAME required=true"`
	User     string `env:"key=DB_USER required=true"`
	Password string `env:"key=DB_PASS required=true"`
	Host     string `env:"key=DB_HOST default=localhost"`
	LogSql   bool   `env:"key=LOG_SQL default=false"`
}

type TestDBConfig struct {
	Name     string `env:"key=TEST_DB_NAME"`
	User     string `env:"key=TEST_DB_USER"`
	Password string `env:"key=TEST_DB_PASS"`
	Host     string `env:"key=TEST_DB_HOST"`
}

type SmtpConfig struct {
	Name          string `env:"key=SMTP_NAME"`
	Password      string `env:"key=SMTP_PASS"`
	User          string `env:"key=SMTP_USER"`
	Host          string `env:"key=SMTP_HOST default=localhost"`
	DefaultSender string `env"key=DEFAULT_EMAIL_SENDER default=webmaster@localhost"`
}

type DirConfig struct {
	Base       string `env:"key=BASE_DIR"`
	Public     string `env:"key=PUBLIC_DIR"`
	Uploads    string `env:"key=UPLOADS_DIR"`
	Thumbnails string `env:"key=THUMBNAILS_DIR"`
	Templates  string `env:"key=TEMPLATES_DIR"`
}

type KeysConfig struct {
	Private string `env:"key=PRIVATE_KEY required=true"`
	Public  string `env:"key=PUBLIC_KEY required=true"`
}

type ServerConfig struct {
	Port int `env:"key=PORT default=5000"`
}

var Config = &struct {
	DB     *DBConfig
	TestDB *TestDBConfig
	Smtp   *SmtpConfig
	Dirs   *DirConfig
	Keys   *KeysConfig
	Server *ServerConfig
}{}

func initConfig() {

	Config.DB = &DBConfig{}

	if err := env.Process(Config.DB); err != nil {
		log.Fatal(err)
	}

	Config.TestDB = &TestDBConfig{}

	if err := env.Process(Config.TestDB); err != nil {
		log.Fatal(err)
	}

	Config.Smtp = &SmtpConfig{}

	if err := env.Process(Config.Smtp); err != nil {
		log.Fatal(err)
	}

	Config.Keys = &KeysConfig{}

	if err := env.Process(Config.Keys); err != nil {
		log.Fatal(err)
	}

	Config.Dirs = &DirConfig{}

	if err := env.Process(Config.Dirs); err != nil {
		log.Fatal(err)
	}

	Config.Server = &ServerConfig{}

	if err := env.Process(Config.Server); err != nil {
		log.Fatal(err)
	}

	if Config.TestDB.Name == "" {
		Config.TestDB.Name = Config.DB.Name + "_test"
	}

	if Config.TestDB.User == "" {
		Config.TestDB.User = Config.DB.User
	}

	if Config.TestDB.Password == "" {
		Config.TestDB.Password = Config.DB.Password
	}

	if Config.TestDB.Host == "" {
		Config.TestDB.Host = Config.DB.Host
	}

	if Config.TestDB.Name == Config.DB.Name {
		log.Fatal("Test DB name same as DB name")
	}

	if Config.Dirs.Base == "" {
		Config.Dirs.Base = getDefaultBaseDir()
	}

	if Config.Dirs.Public == "" {
		Config.Dirs.Public = path.Join(Config.Dirs.Base, "public")
	}

	if Config.Dirs.Uploads == "" {
		Config.Dirs.Uploads = path.Join(Config.Dirs.Public, "uploads")
	}

	if Config.Dirs.Thumbnails == "" {
		Config.Dirs.Thumbnails = path.Join(Config.Dirs.Uploads, "thumbnails")
	}

	if Config.Dirs.Templates == "" {
		Config.Dirs.Templates = path.Join(Config.Dirs.Base, "templates")
	}
}

func getDefaultBaseDir() string {
	defaultBaseDir, err := os.Getwd()
	if err != nil {
		defaultBaseDir = "."
	}
	return defaultBaseDir
}
