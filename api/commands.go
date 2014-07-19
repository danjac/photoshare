package api

import (
	"fmt"
	"github.com/zenazn/goji/graceful"
	"log"
	"runtime"
)

func Serve() {

	config, err := NewAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)

	db, err := InitDB(config)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		log.Fatal("Closing DB connection")
		db.Close()
	}()

	router, err := GetRouter(config, db)
	if err != nil {
		log.Fatal(err)
	}

	if err := graceful.ListenAndServe(fmt.Sprintf(":%d", config.ServerPort), router); err != nil {
		log.Fatal(err)
	}

	graceful.Wait()
}
