package photoshare

import (
	"flag"
	"fmt"
	"github.com/codegangsta/negroni"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Serve runs the HTTP server
func Serve() {

	app, err := newApp()
	if err != nil {
		log.Fatal(err)
	}
	defer app.close()

	runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)

	n := negroni.Classic()
	n.UseHandler(app.getRouter())
	n.Run(fmt.Sprintf(":%d", app.cfg.ServerPort))

}

func storeFile(app *app,
	filename,
	title,
	contentType string,
	tags []string,
	userID int64) error {
	log.Println(title)
	name := generateRandomFilename(contentType)
	file, err := os.Open(filename)
	if err != nil {
		logError(err)
	}
	defer file.Close()
	err = app.filestore.store(file, name, contentType)
	if err != nil {
		logError(err)
	}
	photo := &photo{
		Title:    title,
		Filename: name,
		Tags:     tags,
		OwnerID:  userID,
	}
	if err := app.datamapper.createPhoto(photo); err != nil {
		return err
	}

	return nil
}

func scanDir(app *app, userID int64, baseDir, dirname string) {
	fileList, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Println(err)
	}
	for _, info := range fileList {
		name := info.Name()
		if info.IsDir() {
			scanDir(app, userID, baseDir, filepath.Join(dirname, name))
		} else {
			fullPath := filepath.Join(dirname, name)
			tags := filepath.SplitList(dirname[len(baseDir):])
			ext := strings.ToLower(filepath.Ext(name))
			if ext != ".jpg" && ext != ".png" {
				continue
			}
			title := name[:len(name)-len(ext)]

			var contentType string
			if ext == ".jpg" {
				contentType = "image/jpeg"
			} else {
				contentType = "image/png"
			}

			if err := storeFile(app, fullPath, title, contentType, tags, userID); err != nil {
				log.Println(err)
			}
		}
	}
}

// Import from a given directory. Subdirs will be tags. Title will be filename.
func Import() {

	email := flag.String("user", "", "User email address")
	dirname := flag.String("dir", "", "Directory")

	flag.Parse()

	fmt.Println(*email)
	fmt.Println(*dirname)

	app, err := newApp()
	if err != nil {
		log.Fatal(err)
	}
	defer app.close()

	user, err := app.datamapper.getUserByEmail(*email)
	if err != nil {
		log.Fatal(err)
	}

	scanDir(app, user.ID, *dirname, *dirname)

}
