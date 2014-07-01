package config

import (
	"github.com/danryan/env"
	"log"
	"os"
	"path"
)

var DB = &struct {
	Name     string `env:"key=DB_NAME required=true"`
	User     string `env:"key=DB_USER required=true"`
	Password string `env:"key=DB_PASS required=true"`
	Host     string `env:"key=DB_HOST default=localhost"`
	LogSql   bool   `env:"key=LOG_SQL default=false"`
}{}

var TestDB = &struct {
	Name     string `env:"key=TEST_DB_NAME"`
	User     string `env:"key=TEST_DB_USER"`
	Password string `env:"key=TEST_DB_PASS"`
	Host     string `env:"key=TEST_DB_HOST"`
}{}

var Smtp = &struct {
	Name          string `env:"key=SMTP_NAME"`
	Password      string `env:"key=SMTP_PASS"`
	User          string `env:"key=SMTP_USER"`
	Host          string `env:"key=SMTP_HOST default=localhost"`
	DefaultSender string `env"key=DEFAULT_EMAIL_SENDER default=webmaster@localhost"`
}{}

var Dirs = &struct {
	Base       string `env:"key=BASE_DIR"`
	Public     string `env:"key=PUBLIC_DIR"`
	Uploads    string `env:"key=UPLOADS_DIR"`
	Thumbnails string `env:"key=THUMBNAILS_DIR"`
	Templates  string `env:"key=TEMPLATES_DIR"`
}{}

var Keys = &struct {
	Private string `env:"key=PRIVATE_KEY required=true"`
	Public  string `env:"key=PUBLIC_KEY required=true"`
}{}

var Server = &struct {
	Port int `env:"key=PORT default=5000"`
}{}

func init() {

	if err := env.Process(DB); err != nil {
		log.Fatal(err)
	}

	if err := env.Process(TestDB); err != nil {
		log.Fatal(err)
	}

	if err := env.Process(Keys); err != nil {
		log.Fatal(err)
	}

	if err := env.Process(Dirs); err != nil {
		log.Fatal(err)
	}

	if err := env.Process(Server); err != nil {
		log.Fatal(err)
	}

	if TestDB.Name == "" {
		TestDB.Name = DB.Name + "_test"
	}

	if TestDB.User == "" {
		TestDB.User = DB.User
	}

	if TestDB.Password == "" {
		TestDB.Password = DB.Password
	}

	if TestDB.Name == DB.Name {
		log.Fatal("Test DB name same as DB name")
	}

	defaultBaseDir, err := os.Getwd()
	if err != nil {
		defaultBaseDir = "."
	}

	if Dirs.Base == "" {
		Dirs.Base = defaultBaseDir
	}

	if Dirs.Public == "" {
		Dirs.Public = path.Join(Dirs.Base, "public")
	}

	if Dirs.Uploads == "" {
		Dirs.Uploads = path.Join(Dirs.Public, "uploads")
	}

	if Dirs.Thumbnails == "" {
		Dirs.Thumbnails = path.Join(Dirs.Uploads, "thumbnails")
	}

	if Dirs.Templates == "" {
		Dirs.Templates = path.Join(Dirs.Base, "templates")
	}
}
