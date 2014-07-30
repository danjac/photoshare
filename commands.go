package photoshare

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"log"
	"runtime"
)

// Serve runs the HTTP server
func Serve() {

	cfg, err := newConfigurator()
	if err != nil {
		log.Fatal(err)
	}
	defer cfg.close()

	runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)

	n := negroni.Classic()
	n.UseHandler(cfg.getRouter())
	n.Run(fmt.Sprintf(":%d", cfg.ServerPort))

}
