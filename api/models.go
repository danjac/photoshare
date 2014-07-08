package api

import (
	"bytes"
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/rand"
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

const (
	pageSize               = 12
	recoveryCodeLength     = 30
	recoveryCodeCharacters = "abcdefghijklmnopqrstuvwxyz0123456789"
)

var db *gorp.DbMap

var userMgr = NewUserManager()
var photoMgr = NewPhotoManager()

func InitDB(conn *sql.DB, logSql bool) (*gorp.DbMap, error) {
	db = &gorp.DbMap{Db: conn, Dialect: gorp.PostgresDialect{}}

	if logSql {
		db.TraceOn("[sql]", log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds))
	}

	db.AddTableWithName(User{}, "users").SetKeys(true, "ID")
	db.AddTableWithName(Photo{}, "photos").SetKeys(true, "ID")
	db.AddTableWithName(Tag{}, "tags").SetKeys(true, "ID")

	return db, nil
}

type PhotoManager interface {
	Insert(*Photo) error
	Update(*Photo) error
	Delete(*Photo) error
	Get(int64) (*Photo, bool, error)
	GetDetail(int64, *User) (*PhotoDetail, bool, error)
	GetTagCounts() ([]TagCount, error)
	All(*Page, string) (*PhotoList, error)
	ByOwnerID(*Page, int64) (*PhotoList, error)
	Search(*Page, string) (*PhotoList, error)
	UpdateTags(*Photo) error
}

type PhotoList struct {
	Items       []Photo `json:"photos"`
	Total       int64   `json:"total"`
	CurrentPage int64   `json:"currentPage"`
	NumPages    int64   `json:"numPages"`
}

func NewPhotoList(photos []Photo, total int64, page int64) *PhotoList {
	numPages := int64(math.Ceil(float64(total) / float64(pageSize)))

	return &PhotoList{
		Items:       photos,
		Total:       total,
		CurrentPage: page,
		NumPages:    numPages,
	}
}

type Tag struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type TagCount struct {
	Name      string `db:"name" json:"name"`
	Photo     string `db:"photo" json:"photo"`
	NumPhotos int64  `db:"num_photos" json:"numPhotos"`
}

type Photo struct {
	ID        int64     `db:"id" json:"id"`
	OwnerID   int64     `db:"owner_id" json:"ownerId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Title     string    `db:"title" json:"title"`
	Filename  string    `db:"photo" json:"photo"`
	Tags      []string  `db:"-" json:"tags,omitempty"`
	UpVotes   int64     `db:"up_votes" json:"upVotes"`
	DownVotes int64     `db:"down_votes" json:"downVotes"`
}

func (photo *Photo) PreInsert(s gorp.SqlExecutor) error {
	photo.CreatedAt = time.Now()
	return nil
}

func (photo *Photo) PreDelete(s gorp.SqlExecutor) error {
	go photoCleaner.Clean(photo.Filename)
	return nil
}

func (photo *Photo) CanEdit(user *User) bool {
	if user == nil || !user.IsAuthenticated {
		return false
	}
	return user.IsAdmin || photo.OwnerID == user.ID
}

func (photo *Photo) CanDelete(user *User) bool {
	return photo.CanEdit(user)
}

func (photo *Photo) CanVote(user *User) bool {
	if user == nil || !user.IsAuthenticated {
		return false
	}
	if photo.OwnerID == user.ID {
		return false
	}

	return !user.HasVoted(photo.ID)
}

type Permissions struct {
	Edit   bool `json:"edit"`
	Delete bool `json:"delete"`
	Vote   bool `json:"vote"`
}

type PhotoDetail struct {
	Photo       `db:"-"`
	OwnerName   string       `db:"owner_name" json:"ownerName"`
	Permissions *Permissions `db:"-" json:"perms"`
}

type defaultPhotoManager struct{}

func NewPhotoManager() PhotoManager {
	return &defaultPhotoManager{}
}

func (mgr *defaultPhotoManager) Delete(photo *Photo) error {
	_, err := db.Delete(photo)
	return err
}

func (mgr *defaultPhotoManager) Update(photo *Photo) error {
	_, err := db.Update(photo)
	return err
}

func (mgr *defaultPhotoManager) Insert(photo *Photo) error {
	t, err := db.Begin()
	if err != nil {
		return err
	}
	if err := db.Insert(photo); err != nil {
		return err
	}
	if err := mgr.UpdateTags(photo); err != nil {
		return err
	}
	return t.Commit()
}

func (mgr *defaultPhotoManager) UpdateTags(photo *Photo) error {

	var (
		args    = []string{"$1"}
		params  = []interface{}{interface{}(photo.ID)}
		isEmpty = true
	)
	for i, name := range photo.Tags {
		name = strings.TrimSpace(name)
		if name != "" {
			args = append(args, fmt.Sprintf("$%d", i+2))
			params = append(params, interface{}(strings.ToLower(name)))
			isEmpty = false
		}
	}
	if isEmpty && photo.ID != 0 {
		_, err := db.Exec("DELETE FROM photo_tags WHERE photo_id=$1", photo.ID)
		return err
	}
	_, err := db.Exec(fmt.Sprintf("SELECT add_tags(%s)", strings.Join(args, ",")), params...)
	return err

}

func (mgr *defaultPhotoManager) Get(photoID int64) (*Photo, bool, error) {

	photo := &Photo{}

	if photoID == 0 {
		return photo, false, nil
	}

	obj, err := db.Get(photo, photoID)
	if err != nil {
		return photo, false, err
	}
	if obj == nil {
		return photo, false, nil
	}
	return obj.(*Photo), true, nil
}

func (mgr *defaultPhotoManager) GetDetail(photoID int64, user *User) (*PhotoDetail, bool, error) {

	photo := &PhotoDetail{}

	if photoID == 0 {
		return photo, false, nil
	}

	q := "SELECT p.*, u.name AS owner_name " +
		"FROM photos p JOIN users u ON u.id = p.owner_id " +
		"WHERE p.id=$1"

	if err := db.SelectOne(photo, q, photoID); err != nil {
		if err == sql.ErrNoRows {
			return photo, false, nil
		}
		return photo, false, err
	}

	var tags []Tag

	if _, err := db.Select(&tags,
		"SELECT t.* FROM tags t JOIN photo_tags pt ON pt.tag_id=t.id "+
			"WHERE pt.photo_id=$1", photo.ID); err != nil {
		return photo, false, err
	}
	for _, tag := range tags {
		photo.Tags = append(photo.Tags, tag.Name)
	}

	photo.Permissions = &Permissions{
		photo.CanEdit(user),
		photo.CanDelete(user),
		photo.CanVote(user),
	}
	return photo, true, nil

}

func (mgr *defaultPhotoManager) ByOwnerID(page *Page, ownerID int64) (*PhotoList, error) {

	var (
		photos []Photo
		err    error
		total  int64
	)
	if ownerID == 0 {
		return nil, nil
	}
	if total, err = db.SelectInt("SELECT COUNT(id) FROM photos WHERE owner_id=$1", ownerID); err != nil {
		return nil, err
	}

	if _, err = db.Select(&photos,
		"SELECT * FROM photos WHERE owner_id = $1"+
			"ORDER BY (up_votes - down_votes) DESC, created_at DESC LIMIT $2 OFFSET $3",
		ownerID, page.Size, page.Offset); err != nil {
		return nil, err
	}
	return NewPhotoList(photos, total, page.Index), nil
}

func (mgr *defaultPhotoManager) Search(page *Page, q string) (*PhotoList, error) {

	var (
		clauses []string
		params  []interface{}
		err     error
		photos  []Photo
		total   int64
	)

	if q == "" {
		return nil, nil
	}

	for num, word := range strings.Split(q, " ") {
		word = strings.TrimSpace(word)
		if word == "" || num > 6 {
			break
		}

		num += 1

		if strings.HasPrefix(word, "@") {
			word = word[1:]
			clauses = append(clauses, fmt.Sprintf(
				"SELECT p.* FROM photos p "+
					"INNER JOIN users u ON u.id = p.owner_id  "+
					"WHERE UPPER(u.name::text) = UPPER($%d)", num))
		} else if strings.HasPrefix(word, "#") {
			word = word[1:]
			clauses = append(clauses, fmt.Sprintf(
				"SELECT p.* FROM photos p "+
					"INNER JOIN photo_tags pt ON pt.photo_id = p.id "+
					"INNER JOIN tags t ON pt.tag_id=t.id "+
					"WHERE UPPER(t.name::text) = UPPER($%d)", num))
		} else {
			word = "%" + word + "%"
			clauses = append(clauses, fmt.Sprintf(
				"SELECT DISTINCT p.* FROM photos p "+
					"INNER JOIN users u ON u.id = p.owner_id  "+
					"LEFT JOIN photo_tags pt ON pt.photo_id = p.id "+
					"LEFT JOIN tags t ON pt.tag_id=t.id "+
					"WHERE UPPER(p.title::text) LIKE UPPER($%d) OR "+
					"UPPER(u.name::text) LIKE UPPER($%d) OR t.name LIKE $%d",
				num, num, num))
		}

		params = append(params, interface{}(word))
	}

	clausesSql := strings.Join(clauses, " INTERSECT ")

	countSql := fmt.Sprintf("SELECT COUNT(id) FROM (%s) q", clausesSql)

	if total, err = db.SelectInt(countSql, params...); err != nil {
		return nil, err
	}

	numParams := len(params)

	sql := fmt.Sprintf("SELECT * FROM (%s) q ORDER BY (up_votes - down_votes) DESC, created_at DESC LIMIT $%d OFFSET $%d",
		clausesSql, numParams+1, numParams+2)

	params = append(params, interface{}(page.Size))
	params = append(params, interface{}(page.Offset))

	if _, err = db.Select(&photos, sql, params...); err != nil {
		return nil, err
	}
	return NewPhotoList(photos, total, page.Index), nil
}

func (mgr *defaultPhotoManager) All(page *Page, orderBy string) (*PhotoList, error) {

	var (
		total  int64
		photos []Photo
		err    error
	)
	if orderBy == "votes" {
		orderBy = "(up_votes - down_votes)"
	} else {
		orderBy = "created_at"
	}

	if total, err = db.SelectInt("SELECT COUNT(id) FROM photos"); err != nil {
		return nil, err
	}

	if _, err = db.Select(&photos,
		"SELECT * FROM photos "+
			"ORDER BY "+orderBy+" DESC LIMIT $1 OFFSET $2", page.Size, page.Offset); err != nil {
		return nil, err
	}
	return NewPhotoList(photos, total, page.Index), nil
}

func (mgr *defaultPhotoManager) GetTagCounts() ([]TagCount, error) {
	var tags []TagCount
	if _, err := db.Select(&tags, "SELECT name, photo, num_photos FROM tag_counts"); err != nil {
		return tags, err
	}
	return tags, nil
}

type UserManager interface {
	Insert(user *User) error
	Update(user *User) error
	IsNameAvailable(user *User) (bool, error)
	IsEmailAvailable(user *User) (bool, error)
	GetActive(userID int64) (*User, bool, error)
	GetByRecoveryCode(string) (*User, bool, error)
	GetByEmail(string) (*User, bool, error)
	Authenticate(identifier string, password string) (*User, bool, error)
}

type defaultUserManager struct{}

func (mgr *defaultUserManager) Insert(user *User) error {
	return db.Insert(user)
}

func (mgr *defaultUserManager) Update(user *User) error {
	_, err := db.Update(user)
	return err
}

func (mgr *defaultUserManager) IsNameAvailable(user *User) (bool, error) {
	var (
		num int64
		err error
	)
	q := "SELECT COUNT(id) FROM users WHERE name=$1"
	if user.ID == 0 {
		num, err = db.SelectInt(q, user.Name)
	} else {
		q += " AND id != $2"
		num, err = db.SelectInt(q, user.Name, user.ID)
	}
	if err != nil {
		return false, err
	}
	return num == 0, nil
}

func (mgr *defaultUserManager) IsEmailAvailable(user *User) (bool, error) {
	var (
		num int64
		err error
	)
	q := "SELECT COUNT(id) FROM users WHERE email=$1"
	if user.ID == 0 {
		num, err = db.SelectInt(q, user.Email)
	} else {
		q += " AND id != $2"
		num, err = db.SelectInt(q, user.Email, user.ID)
	}
	if err != nil {
		return false, err
	}
	return num == 0, nil
}
func (mgr *defaultUserManager) GetActive(userID int64) (*User, bool, error) {

	user := &User{}
	if err := db.SelectOne(user, "SELECT * FROM users WHERE active=$1 AND id=$2", true, userID); err != nil {
		if err == sql.ErrNoRows {
			return user, false, nil
		}
		return user, false, err
	}
	return user, true, nil

}

func (mgr *defaultUserManager) GetByRecoveryCode(code string) (*User, bool, error) {

	user := &User{}
	if code == "" {
		return user, false, nil
	}
	if err := db.SelectOne(user, "SELECT * FROM users WHERE active=$1 AND recovery_code=$2", true, code); err != nil {
		if err == sql.ErrNoRows {
			return user, false, nil
		}
		return user, false, err
	}
	return user, true, nil

}
func (mgr *defaultUserManager) GetByEmail(email string) (*User, bool, error) {
	user := &User{}
	if err := db.SelectOne(user, "SELECT * FROM users WHERE active=$1 AND email=$2", true, email); err != nil {
		if err == sql.ErrNoRows {
			return user, false, nil
		}
		return user, false, err
	}
	return user, true, nil
}

func (mgr *defaultUserManager) Authenticate(identifier, password string) (*User, bool, error) {
	user := &User{}

	if err := db.SelectOne(user, "SELECT * FROM users WHERE active=$1 AND (email=$2 OR name=$2)", true, identifier); err != nil {
		if err == sql.ErrNoRows {
			return user, false, nil
		}
		return user, false, err
	}

	if !user.CheckPassword(password) {
		return user, false, nil
	}

	return user, true, nil
}

func NewUserManager() UserManager {
	return &defaultUserManager{}
}

type User struct {
	ID              int64          `db:"id" json:"id"`
	CreatedAt       time.Time      `db:"created_at" json:"createdAt"`
	Name            string         `db:"name" json:"name"`
	Password        string         `db:"password" json:""`
	Email           string         `db:"email" json:"email"`
	Votes           string         `db:"votes" json:""`
	IsAdmin         bool           `db:"admin" json:"isAdmin"`
	IsActive        bool           `db:"active" json:"isActive"`
	RecoveryCode    sql.NullString `db:"recovery_code" json:""`
	IsAuthenticated bool           `db:"-" json:"isAuthenticated"`
}

func (user *User) PreInsert(s gorp.SqlExecutor) error {
	user.IsActive = true
	user.CreatedAt = time.Now()
	user.EncryptPassword()
	user.Votes = "{}"
	return nil
}

func (user *User) GenerateRecoveryCode() (string, error) {

	buf := bytes.Buffer{}
	randbytes := make([]byte, recoveryCodeLength)

	if _, err := rand.Read(randbytes); err != nil {
		return "", err
	}

	numChars := len(recoveryCodeCharacters)

	for i := 0; i < recoveryCodeLength; i++ {
		index := int(randbytes[i]) % numChars
		char := recoveryCodeCharacters[index]
		buf.WriteString(string(char))
	}

	code := buf.String()
	user.RecoveryCode = sql.NullString{String: code, Valid: true}
	return code, nil
}

func (user *User) ResetRecoveryCode() {
	user.RecoveryCode = sql.NullString{String: "", Valid: false}
}

func (user *User) ChangePassword(password string) error {
	user.Password = password
	return user.EncryptPassword()
}

func (user *User) EncryptPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
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

func (user *User) RegisterVote(photoID int64) {
	user.SetVotes(append(user.GetVotes(), photoID))
}

func (user *User) HasVoted(photoID int64) bool {
	for _, value := range user.GetVotes() {
		if value == photoID {
			return true
		}
	}
	return false
}
func (user *User) GetVotes() []int64 {
	return pgArrToIntSlice(user.Votes)
}

func (user *User) SetVotes(votes []int64) {
	user.Votes = intSliceToPgArr(votes)
}

type Page struct {
	Index  int64
	Offset int64
	Size   int64
}

func NewPage(index int64) *Page {
	offset := (index - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	return &Page{index, offset, pageSize}
}
