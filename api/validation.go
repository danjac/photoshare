package api

import (
	"regexp"
)

var (
	emailRegex  = regexp.MustCompile(".+@.+\\..+")
	formHandler = &FormHandler{}
)

type FormHandler struct{}

func (h *FormHandler) Validate(validator Validator) error {
	result := NewValidationResult()
	if err := validator.Validate(result); err != nil {
		return err
	}
	if len(result.Errors) > 0 {
		return result
	}
	return nil
}

type Validator interface {
	Validate(result *ValidationResult) error
}

type ValidationResult struct {
	Errors map[string]string `json:"errors"`
}

func (result ValidationResult) Error() string {
	return "Validation errors"
}

var getPhotoValidator = func(photo *Photo) Validator {
	return NewPhotoValidator(photo)
}

var getUserValidator = func(user *User) Validator {
	return NewUserValidator(user)
}

func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		make(map[string]string),
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
		result.Errors["ownerID"] = "Owner ID is missing"
	}
	if v.Photo.Title == "" {
		result.Errors["title"] = "Title is missing"
	}
	if len(v.Photo.Title) > 200 {
		result.Errors["title"] = "Title is too long"
	}
	if v.Photo.Filename == "" {
		result.Errors["photo"] = "Photo filename not set"
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
		result.Errors["name"] = "Name is missing"
	} else {
		ok, err := userMgr.IsNameAvailable(v.User)
		if err != nil {
			return err
		}
		if !ok {
			result.Errors["name"] = "Name already taken"
		}
	}

	if v.User.Email == "" {
		result.Errors["email"] = "Email is missing"
	} else if !validateEmail(v.User.Email) {
		result.Errors["email"] = "Invalid email address"
	} else {
		ok, err := userMgr.IsEmailAvailable(v.User)
		if err != nil {
			return err
		}
		if !ok {
			result.Errors["email"] = "Email already taken"
		}

	}

	if v.User.Password == "" {
		result.Errors["password"] = "Password is missing"
	}

	return nil

}
