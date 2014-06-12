package validation

import (
	"github.com/danjac/photoshare/api/models"
	"regexp"
)

var (
	userMgr    = models.NewUserManager()
	emailRegex = regexp.MustCompile(".+@.+\\..+")
)

func validateEmail(email string) bool {
	return emailRegex.Match([]byte(email))
}

type UserValidator struct {
	User *models.User
}

func (v *UserValidator) Validate() (*ValidationResult, error) {

	result := NewValidationResult()

	if v.User.Name == "" {
		result.Error("name", "Name is missing")
	} else {
		ok, err := userMgr.IsNameAvailable(v.User)
		if err != nil {
			return result, err
		}
		if !ok {
			result.Error("name", "Name already taken")
		}
	}

	if v.User.Email == "" {
		result.Error("email", "Email is missing")
	} else if !validateEmail(v.User.Email) {
		result.Error("email", "Invalid email address")
	} else {
		ok, err := userMgr.IsEmailAvailable(v.User)
		if err != nil {
			return result, err
		}
		if !ok {
			result.Error("email", "Email already taken")
		}

	}

	if v.User.Password == "" {
		result.Error("password", "Password is missing")
	}

	return result, nil

}
