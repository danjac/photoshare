package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"github.com/coopernurse/gorp"
	"time"
)

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

func (user *User) Save() error {
	return dbMap.Insert(user)
}

func NewUser(name, email, password string) *User {
	user := &User{Name: name, Email: email, IsActive: true}
	user.SetPassword(password)
	return user
}

func GetUser(userID int) (*User, error) {
	obj, err := dbMap.Get(User{}, userID)
	if err != nil {
		return nil, err
	}
	if obj == nil {
		return nil, nil
	}

	return obj.(*User), nil
}

func Authenticate(email string, password string) (*User, error) {
	user := &User{}
	if err := dbMap.SelectOne(user, "SELECT * FROM users WHERE active=1 AND email=?", email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if !user.CheckPassword(password) {
		return nil, nil
	}

	return user, nil
}
