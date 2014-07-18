package api

import (
	"bytes"
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/rand"
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"math"
	"time"
)

const (
	pageSize               = 12
	recoveryCodeLength     = 30
	recoveryCodeCharacters = "abcdefghijklmnopqrstuvwxyz0123456789"
)

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

/*
func (photo *Photo) PreDelete(s gorp.SqlExecutor) error {
	go fileMgr.Clean(photo.Filename)
	return nil
}
*/

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
