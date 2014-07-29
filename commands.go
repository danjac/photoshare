package photoshare

import (
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
	defer config.close()

	runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)

	n := negroni.Classic()
	n.UseHandler(config.getRouter())
	n.Run(fmt.Sprintf(":%d", config.ServerPort))

}
