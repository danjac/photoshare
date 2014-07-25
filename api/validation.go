package api

import (
	"regexp"
)

var emailRegex = regexp.MustCompile(".+@.+\\..+")

// ValidationFailure represents a set of validation errors
type ValidationFailure struct {
	Errors map[string]string `json:"errors"`
}

func (err ValidationFailure) Error() string {
	return "Validation failure"
}

func validate(validator Validator) error {
	errors := make(map[string]string)
	if err := validator.Validate(errors); err != nil {
		return err
	}
	if len(errors) > 0 {
		return ValidationFailure{errors}
	}
	return nil
}

// Validator is a common validation interface
type Validator interface {
	Validate(map[string]string) error
}

// NewPhotoValidator creates a new PhotoValidator instance
func NewPhotoValidator(photo *Photo) *PhotoValidator {
	return &PhotoValidator{photo}
}

// PhotoValidator checks if a photo is valid
type PhotoValidator struct {
	photo *Photo
}

// Validate does actual validation
func (v *PhotoValidator) Validate(errors map[string]string) error {
	if v.photo.OwnerID == 0 {
		errors["ownerID"] = "Owner ID is missing"
	}
	if v.photo.Title == "" {
		errors["title"] = "Title is missing"
	}
	if len(v.photo.Title) > 200 {
		errors["title"] = "Title is too long"
	}
	if v.photo.Filename == "" {
		errors["photo"] = "Photo filename not set"
	}
	return nil
}

func validateEmail(email string) bool {
	return emailRegex.Match([]byte(email))
}

// NewUserValidator creates new UserValidator instance
func NewUserValidator(user *User, mgr UserDataStore) *UserValidator {
	return &UserValidator{user, mgr}
}

// UserValidator validates user model is correct
type UserValidator struct {
	user   *User
	userDS UserDataStore
}

// Validate does actual validation
func (v *UserValidator) Validate(errors map[string]string) error {

	if v.user.Name == "" {
		errors["name"] = "Name is missing"
	} else {
		ok, err := v.userDS.IsNameAvailable(v.user)
		if err != nil {
			return err
		}
		if !ok {
			errors["name"] = "Name already taken"
		}
	}

	if v.user.Email == "" {
		errors["email"] = "Email is missing"
	} else if !validateEmail(v.user.Email) {
		errors["email"] = "Invalid email address"
	} else {
		ok, err := v.userDS.IsEmailAvailable(v.user)
		if err != nil {
			return err
		}
		if !ok {
			errors["email"] = "Email already taken"
		}

	}

	if v.user.Password == "" {
		errors["password"] = "Password is missing"
	}

	return nil

}
