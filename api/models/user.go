package models

import (
	"bytes"
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/rand"
	"database/sql"
	"github.com/coopernurse/gorp"
	"time"
)

const (
	recoveryCodeLength     = 30
	recoveryCodeCharacters = "abcdefghijklmnopqrstuvwxyz0123456789"
)

type UserManager interface {
	Insert(user *User) error
	Update(user *User) error
	IsNameAvailable(user *User) (bool, error)
	IsEmailAvailable(user *User) (bool, error)
	GetActive(userID int64) (*User, error)
	GetByRecoveryCode(string) (*User, error)
	GetByEmail(string) (*User, error)
	Authenticate(identifier string, password string) (*User, error)
}

type defaultUserManager struct{}

func (mgr *defaultUserManager) Insert(user *User) error {
	return dbMap.Insert(user)
}

func (mgr *defaultUserManager) Update(user *User) error {
	_, err := dbMap.Update(user)
	return err
}

func (mgr *defaultUserManager) IsNameAvailable(user *User) (bool, error) {
	var (
		num int64
		err error
	)
	q := "SELECT COUNT(id) FROM users WHERE name=$1"
	if user.ID == 0 {
		num, err = dbMap.SelectInt(q, user.Name)
	} else {
		q += " AND id != $2"
		num, err = dbMap.SelectInt(q, user.Name, user.ID)
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
		num, err = dbMap.SelectInt(q, user.Email)
	} else {
		q += " AND id != $2"
		num, err = dbMap.SelectInt(q, user.Email, user.ID)
	}
	if err != nil {
		return false, err
	}
	return num == 0, nil
}
func (mgr *defaultUserManager) GetActive(userID int64) (*User, error) {

	user := &User{}
	if err := dbMap.SelectOne(user, "SELECT * FROM users WHERE active=$1 AND id=$2", true, userID); err != nil {
		if err == sql.ErrNoRows {
			return user, nil
		}
		return user, err
	}
	user.IsAuthenticated = true
	return user, nil

}

func (mgr *defaultUserManager) GetByRecoveryCode(code string) (*User, error) {

	user := &User{}
	if code == "" {
		return user, nil
	}
	if err := dbMap.SelectOne(user, "SELECT * FROM users WHERE active=$1 AND recovery_code=$2", true, code); err != nil {
		if err == sql.ErrNoRows {
			return user, nil
		}
		return user, err
	}
	user.IsAuthenticated = true
	return user, nil

}
func (mgr *defaultUserManager) GetByEmail(email string) (*User, error) {
	user := &User{}
	if err := dbMap.SelectOne(user, "SELECT * FROM users WHERE active=$1 AND email=$2", true, email); err != nil {
		if err == sql.ErrNoRows {
			return user, nil
		}
		return user, err
	}
	user.IsAuthenticated = true
	return user, nil
}

func (mgr *defaultUserManager) Authenticate(identifier, password string) (*User, error) {
	user := &User{}
	if err := dbMap.SelectOne(user, "SELECT * FROM users WHERE active=$1 AND (email=$2 OR name=$2)", true, identifier); err != nil {
		if err == sql.ErrNoRows {
			return user, nil
		}
		return user, err
	}

	if !user.CheckPassword(password) {
		return user, nil
	}

	user.IsAuthenticated = true

	return user, nil
}

var userMgr = &defaultUserManager{}

func NewUserManager() UserManager {
	return userMgr
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
