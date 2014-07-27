package photoshare

import (
	"bytes"
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/rand"
	"database/sql"
	"github.com/coopernurse/gorp"
	"math"
	"time"
)

const (
	pageSize               = 20
	recoveryCodeLength     = 30
	recoveryCodeCharacters = "abcdefghijklmnopqrstuvwxyz0123456789"
)

type photoList struct {
	Items       []photo `json:"photos"`
	Total       int64   `json:"total"`
	CurrentPage int64   `json:"currentPage"`
	NumPages    int64   `json:"numPages"`
}

func newPhotoList(photos []photo, total int64, page int64) *photoList {
	numPages := int64(math.Ceil(float64(total) / float64(pageSize)))

	return &photoList{
		Items:       photos,
		Total:       total,
		CurrentPage: page,
		NumPages:    numPages,
	}
}

type tag struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type tagCount struct {
	Name      string `db:"name" json:"name"`
	Photo     string `db:"photo" json:"photo"`
	NumPhotos int64  `db:"num_photos" json:"numPhotos"`
}

type photo struct {
	ID        int64     `db:"id" json:"id"`
	OwnerID   int64     `db:"owner_id" json:"ownerId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Title     string    `db:"title" json:"title"`
	Filename  string    `db:"photo" json:"photo"`
	Tags      []string  `db:"-" json:"tags,omitempty"`
	UpVotes   int64     `db:"up_votes" json:"upVotes"`
	DownVotes int64     `db:"down_votes" json:"downVotes"`
}

func (photo *photo) PreInsert(s gorp.SqlExecutor) error {
	photo.CreatedAt = time.Now()
	return nil
}

func (photo *photo) validate(c *appContext, errors map[string]string) error {
	if photo.OwnerID == 0 {
		errors["ownerID"] = "Owner ID is missing"
	}
	if photo.Title == "" {
		errors["title"] = "Title is missing"
	}
	if len(photo.Title) > 200 {
		errors["title"] = "Title is too long"
	}
	if photo.Filename == "" {
		errors["photo"] = "Photo filename not set"
	}
	return nil
}

func (photo *photo) canEdit(user *user) bool {
	if user == nil || !user.IsAuthenticated {
		return false
	}
	return user.IsAdmin || photo.OwnerID == user.ID
}

func (photo *photo) canDelete(user *user) bool {
	return photo.canEdit(user)
}

func (photo *photo) canVote(user *user) bool {
	if user == nil || !user.IsAuthenticated {
		return false
	}
	if photo.OwnerID == user.ID {
		return false
	}

	return !user.hasVoted(photo.ID)
}

type permissions struct {
	Edit   bool `json:"edit"`
	Delete bool `json:"delete"`
	Vote   bool `json:"vote"`
}

type photoDetail struct {
	photo       `db:"-"`
	OwnerName   string       `db:"owner_name" json:"ownerName"`
	Permissions *permissions `db:"-" json:"perms"`
}

// User represents users in database
type user struct {
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

// PreInsert hook
func (user *user) PreInsert(s gorp.SqlExecutor) error {
	user.IsActive = true
	user.CreatedAt = time.Now()
	user.Votes = "{}"
	user.encryptPassword()
	return nil
}

func (user *user) validate(c *appContext, errors map[string]string) error {

	if user.Name == "" {
		errors["name"] = "Name is missing"
	} else {
		ok, err := c.ds.isUserNameAvailable(user)
		if err != nil {
			return err
		}
		if !ok {
			errors["name"] = "Name already taken"
		}
	}

	if user.Email == "" {
		errors["email"] = "Email is missing"
	} else if !validateEmail(user.Email) {
		errors["email"] = "Invalid email address"
	} else {
		ok, err := c.ds.isUserEmailAvailable(user)
		if err != nil {
			return err
		}
		if !ok {
			errors["email"] = "Email already taken"
		}

	}

	if user.Password == "" {
		errors["password"] = "Password is missing"
	}

	return nil

}
func (user *user) generateRecoveryCode() (string, error) {

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

func (user *user) resetRecoveryCode() {
	user.RecoveryCode = sql.NullString{String: "", Valid: false}
}

func (user *user) changePassword(password string) error {
	user.Password = password
	return user.encryptPassword()
}

func (user *user) encryptPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return nil
}

func (user *user) checkPassword(password string) bool {
	if user.Password == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (user *user) registerVote(photoID int64) {
	user.setVotes(append(user.getVotes(), photoID))
}

func (user *user) hasVoted(photoID int64) bool {
	for _, value := range user.getVotes() {
		if value == photoID {
			return true
		}
	}
	return false
}

func (user *user) getVotes() []int64 {
	return pgArrToIntSlice(user.Votes)
}

func (user *user) setVotes(votes []int64) {
	user.Votes = intSliceToPgArr(votes)
}

type page struct {
	index  int64
	offset int64
	size   int64
}

func newPage(index int64) *page {
	offset := (index - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	return &page{index, offset, pageSize}
}
