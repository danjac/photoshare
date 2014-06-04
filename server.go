package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/gorilla/securecookie"
	_ "github.com/mattn/go-sqlite3"
    "github.com/dchest/uniuri"
	"io"
	"net/http"
	"os"
	"time"
)

const CookieName = "userid"

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
	if err := dbMap.TruncateTables(); err != nil {
		panic(err)
	}

	user := &User{Name: "demo", Email: "demo@photoshare.com", IsActive: true}
	user.SetPassword("demo1")
	if err := dbMap.Insert(user); err != nil {
		panic(err)
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

func addPhoto(w http.ResponseWriter, r *http.Request) {

	user, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("no user found"))
		return
	}

	title := r.FormValue("title")
	src, hdr, err := r.FormFile("photo")
    contentType := hdr.Header["Content-Type"][0]
    var ext string

    if contentType == "image/png" {
        ext = ".png"
    } else {
        ext = ".jpg" 
    }

    filename := uniuri.New() + ext

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer src.Close()

	dst, err := os.Create(fmt.Sprintf("app/uploads/%s", filename))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	photo := &Photo{Title: title,
		Photo:   filename,
		OwnerID: user.ID}
	if err := dbMap.Insert(photo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(photo)
}

func getPhotos(w http.ResponseWriter, r *http.Request) {

	var photos []Photo
	if _, err := dbMap.Select(&photos, "SELECT * FROM photos ORDER BY created_at DESC"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(photos)
}

// this should be DELETE
func logout(w http.ResponseWriter, r *http.Request) {

	encoded, err := sCookie.Encode(CookieName, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  CookieName,
		Value: encoded,
		Path:  "/",
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out"))

}

// return current logged in user, or 401
func authenticate(w http.ResponseWriter, r *http.Request) {

	user, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("no user found"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func login(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email or password empty"))
		return
	}

	user := &User{}
	if err := dbMap.SelectOne(user, "SELECT * FROM users WHERE active=1 AND email=?", email); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("User not found"))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !user.CheckPassword(password) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid password"))
		return
	}

	if _, err := json.Marshal(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write the user ID to the secure cookie
	encoded, err := sCookie.Encode(CookieName, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(encoded)

	cookie := &http.Cookie{
		Name:  CookieName,
		Value: encoded,
		Path:  "/",
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(user)

}

func main() {
	db, err := sql.Open("sqlite3", "photos.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	go populateDatabase(db)
	fmt.Println("starting server...")

	// STATIC FILES

	http.HandleFunc("/auth", authenticate)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/photos", getPhotos)
	http.HandleFunc("/add", addPhoto)
	http.Handle("/", http.FileServer(http.Dir("./app/")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	http.ListenAndServe(":"+port, nil)
}
