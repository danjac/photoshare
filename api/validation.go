package api

import (
	"net/http"
	"regexp"
)

var (
	emailRegex  = regexp.MustCompile(".+@.+\\..+")
	formHandler = &FormHandler{}
)

type FormHandler struct{}

func (h *FormHandler) Validate(validator Validator) (*ValidationResult, error) {
	result := NewValidationResult()
	err := validator.Validate(result)
	return result, err
}

type Validator interface {
	Validate(result *ValidationResult) error
}

type ValidationResult struct {
	Errors map[string]string `json:"errors"`
	OK     bool              `json:"ok"`
}

func (result *ValidationResult) Write(w http.ResponseWriter) {
	writeJSON(w, result, http.StatusBadRequest)
}

var getPhotoValidator = func(photo *Photo) Validator {
	return NewPhotoValidator(photo)
}

var getUserValidator = func(user *User) Validator {
	return NewUserValidator(user)
}

func (result *ValidationResult) Error(name, msg string) {
	result.Errors[name] = msg
	result.OK = false
}

func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		make(map[string]string),
		true,
	}
}

func NewPhotoValidator(photo *Photo) *PhotoValidator {
	return &PhotoValidator{photo}
}

type PhotoValidator struct {
	Photo *Photo
}

func (v *PhotoValidator) Validate(result *ValidationResult) error {
	if v.Photo.OwnerID == 0 {
		result.Error("ownerID", "Owner ID is missing")
	}
	if v.Photo.Title == "" {
		result.Error("title", "Title is missing")
	}
	if len(v.Photo.Title) > 200 {
		result.Error("title", "Title is too long")
	}
	if v.Photo.Filename == "" {
		result.Error("photo", "Photo filename not set")
	}
	return nil
}

func validateEmail(email string) bool {
	return emailRegex.Match([]byte(email))
}

func NewUserValidator(user *User) *UserValidator {
	return &UserValidator{user}
}

type UserValidator struct {
	User *User
}

func (v *UserValidator) Validate(result *ValidationResult) error {

	if v.User.Name == "" {
		result.Error("name", "Name is missing")
	} else {
		ok, err := userMgr.IsNameAvailable(v.User)
		if err != nil {
			return err
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
			return err
		}
		if !ok {
			result.Error("email", "Email already taken")
		}

	}

	if v.User.Password == "" {
		result.Error("password", "Password is missing")
	}

	return nil

}
