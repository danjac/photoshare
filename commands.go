package photoshare

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/negroni"
	"log"
	"runtime"
)

// Serve runs the HTTP server
func Serve() {

	config, err := newAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s",
		config.DBUser,
		config.DBName,
		config.DBPassword,
		config.DBHost,
	))

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		log.Fatal("Closing DB connection")
		db.Close()
	}()

	dbMap, err := initDB(db, config.LogSql)
	if err != nil {
		log.Fatal(err)
	}

	context, err := newAppContext(config, dbMap)
	if err != nil {
		log.Fatal(err)
	}

	router, err := getRouter(config, context)
	if err != nil {
		log.Fatal(err)
	}

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(fmt.Sprintf(":%d", config.ServerPort))

}
