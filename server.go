package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	CookieName = "userid"
	UploadsDir = "app/uploads"
)

var dbMap *gorp.DbMap
var hashKey = securecookie.GenerateRandomKey(32)
var blockKey = securecookie.GenerateRandomKey(32)
var sCookie = securecookie.New(hashKey, blockKey)

type Photo struct {
	ID        int       `db:"id" json:"id"`
	OwnerID   int       `db:"owner_id" json:"ownerId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Title     string    `db:"title" json:"title"`
	Photo     string    `db:"photo" json:"photo"`
	Thumbnail string    `db:"thumbnail" json:"thumbnail"`
}

func (photo *Photo) PreInsert(s gorp.SqlExecutor) error {
	photo.CreatedAt = time.Now()
	return nil
}

func (photo *Photo) processImage(src multipart.File, filename, contentType string) error {
	if err := os.MkdirAll(UploadsDir+"/thumbnails", 0777); err != nil && !os.IsExist(err) {
		return err
	}

	// make thumbnail
	var (
		img image.Image
		err error
	)

	if contentType == "image/png" {
		img, err = png.Decode(src)
	} else {
		img, err = jpeg.Decode(src)
	}

	if err != nil {
		return err
	}

	thumb := resize.Thumbnail(300, 300, img, resize.Lanczos3)
	dst, err := os.Create(strings.Join([]string{UploadsDir, "thumbnails", filename}, "/"))

	if err != nil {
		return err
	}

	defer dst.Close()

	if contentType == "image/png" {
		png.Encode(dst, thumb)
	} else if contentType == "image/jpeg" {
		jpeg.Encode(dst, thumb, nil)
	}

	src.Seek(0, 0)

	dst, err = os.Create(strings.Join([]string{UploadsDir, filename}, "/"))

	if err != nil {
		return err
	}

	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

    photo.Photo = filename
    if _, err := dbMap.Update(photo); err != nil {
        return err
    }

	return nil
}


type User struct {
	ID        int       `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Name      string    `db:"name" json:"name"`
	Password  string    `db:"password" json:"-"`
	Email     string    `db:"email" json:"email"`
	IsAdmin   bool      `db:"admin" json:"isAdmin"`
	IsActive  bool      `db:"active" json:"isActive"`
}

func (user *User) PreInsert(s gorp.SqlExecutor) error {
	user.CreatedAt = time.Now()
	return nil
}

func (user *User) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return nil
}

func (user *User) CheckPassword(password string) bool {
	if user.Password == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func populateDatabase(db *sql.DB) {

	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.AddTableWithName(User{}, "users").SetKeys(true, "ID")
	dbMap.AddTableWithName(Photo{}, "photos").SetKeys(true, "ID")
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		panic(err)
	}

	numUsers, err := dbMap.SelectInt("SELECT COUNT(id) FROM users")
	if err != nil {
		panic(err)
	} else if numUsers == 0 {
		user := &User{Name: "demo", Email: "demo@photoshare.com", IsActive: true}
		user.SetPassword("demo1")
		if err := dbMap.Insert(user); err != nil {
			panic(err)
		}
	}

	fmt.Println("Database ready!")

}

func getCurrentUser(r *http.Request) (*User, error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return nil, nil
	}

	var userID int
	if err := sCookie.Decode(CookieName, cookie.Value, &userID); err != nil {
		return nil, nil
	}

	if userID == 0 {
		return nil, nil
	}

	obj, err := dbMap.Get(User{}, userID)
	if err != nil {
		return nil, err
	}
	if obj == nil {
		return nil, nil
	}

	return obj.(*User), nil
}

func renderError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func renderStatus(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

func renderJSON(w http.ResponseWriter, status int, value interface{}) {
	w.WriteHeader(status)
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(value)
}

func upload(w http.ResponseWriter, r *http.Request) {

	user, err := getCurrentUser(r)
	if err != nil {
		renderError(w, err)
		return
	}
	if user == nil {
		renderStatus(w, http.StatusUnauthorized, "Not logged in")
		return
	}

	title := r.FormValue("title")
	src, hdr, err := r.FormFile("photo")
    if err != nil {
        renderError(w, err)
        return
    }
	contentType := hdr.Header["Content-Type"][0]

	defer src.Close()

	filename := uniuri.New()

	if contentType == "image/png" {
		filename += ".png"
	} else if contentType == "image/jpeg" {
		filename += ".jpg"
	} else {
		renderStatus(w, http.StatusBadRequest, "Not a valid image")
		return
	}

	photo := &Photo{Title: title,
		Photo:   filename,
		OwnerID: user.ID}
	if err := dbMap.Insert(photo); err != nil {
		renderError(w, err)
		return
	}

    go photo.processImage(src, filename, contentType)

	renderJSON(w, http.StatusOK, photo)
}

func getPhotos(w http.ResponseWriter, r *http.Request) {

	var photos []Photo
	if _, err := dbMap.Select(&photos, "SELECT * FROM photos WHERE photo != '' AND photo IS NOT NULL  ORDER BY created_at DESC"); err != nil {
		renderError(w, err)
		return
	}
	renderJSON(w, http.StatusOK, photos)
}

// this should be DELETE
func logout(w http.ResponseWriter, r *http.Request) {

	if err := writeCookie(w, 0); err != nil {
		renderError(w, err)
		return
	}

	renderStatus(w, http.StatusOK, "Logged out")

}

func signup(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		renderStatus(w, http.StatusBadRequest, "Missing form info")
		return
	}

	numUsers, err := dbMap.SelectInt("SELECT COUNT(id) FROM users WHERE email=?", email)
	if err != nil {
		renderError(w, err)
		return
	}
	if numUsers > 0 {
		renderStatus(w, http.StatusBadRequest, "Email already taken")
		return
	}

	user := &User{Email: email}
	user.SetPassword(password)
	if err := dbMap.Insert(user); err != nil {
		renderError(w, err)
		return
	}
	if err := writeCookie(w, user.ID); err != nil {
		renderError(w, err)
		return
	}
	renderJSON(w, http.StatusOK, user)
}

// return current logged in user, or 401
func authenticate(w http.ResponseWriter, r *http.Request) {

	user, err := getCurrentUser(r)
	if err != nil {
		renderError(w, err)
		return
	}

	var status int

	if user != nil {
		status = http.StatusOK
	} else {
		status = http.StatusNotFound
	}

	renderJSON(w, status, user)
}

func login(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		renderStatus(w, http.StatusBadRequest, "Email or password missing")
		return
	}

	user := &User{}
	if err := dbMap.SelectOne(user, "SELECT * FROM users WHERE active=1 AND email=?", email); err != nil {
		if err == sql.ErrNoRows {
			renderStatus(w, http.StatusNotFound, "No user found")
			return
		}
		renderError(w, err)
		return
	}

	if !user.CheckPassword(password) {
		renderStatus(w, http.StatusBadRequest, "Invalid password")
		return
	}

	if _, err := json.Marshal(user); err != nil {
		renderError(w, err)
		return
	}

	if err := writeCookie(w, user.ID); err != nil {
		renderError(w, err)
		return
	}

	renderJSON(w, http.StatusOK, user)
}

func writeCookie(w http.ResponseWriter, userID int) error {

	// write the user ID to the secure cookie
	encoded, err := sCookie.Encode(CookieName, userID)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:  CookieName,
		Value: encoded,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	return nil

}

func main() {
	db, err := sql.Open("sqlite3", "photos.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	go populateDatabase(db)

	r := mux.NewRouter()

	s := r.PathPrefix("/auth").Subrouter()
	s.HandleFunc("/", authenticate).Methods("GET")
	s.HandleFunc("/", login).Methods("POST")
	s.HandleFunc("/", logout).Methods("DELETE")

	s = r.PathPrefix("/photos").Subrouter()
	s.HandleFunc("/", getPhotos).Methods("GET")
	s.HandleFunc("/", upload).Methods("POST")

	s = r.PathPrefix("/user").Subrouter()
	s.HandleFunc("/", signup).Methods("POST")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./app/")))

	http.Handle("/", r)

	fmt.Println("starting server...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	http.ListenAndServe(":"+port, nil)
}
