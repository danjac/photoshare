package api

import (
	"regexp"
)

var emailRegex = regexp.MustCompile(".+@.+\\..+")

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

type Validator interface {
	Validate(map[string]string) error
}

var getPhotoValidator = func(photo *Photo) Validator {
	return NewPhotoValidator(photo)
}

var getUserValidator = func(user *User) Validator {
	return NewUserValidator(user)
}

func NewPhotoValidator(photo *Photo) *PhotoValidator {
	return &PhotoValidator{photo}
}

type PhotoValidator struct {
	Photo *Photo
}

func (v *PhotoValidator) Validate(errors map[string]string) error {
	if v.Photo.OwnerID == 0 {
		errors["ownerID"] = "Owner ID is missing"
	}
	if v.Photo.Title == "" {
		errors["title"] = "Title is missing"
	}
	if len(v.Photo.Title) > 200 {
		errors["title"] = "Title is too long"
	}
	if v.Photo.Filename == "" {
		errors["photo"] = "Photo filename not set"
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

func (v *UserValidator) Validate(errors map[string]string) error {

	if v.User.Name == "" {
		errors["name"] = "Name is missing"
	} else {
		ok, err := userMgr.IsNameAvailable(v.User)
		if err != nil {
			return err
		}
		if !ok {
			errors["name"] = "Name already taken"
		}
	}

	if v.User.Email == "" {
		errors["email"] = "Email is missing"
	} else if !validateEmail(v.User.Email) {
		errors["email"] = "Invalid email address"
	} else {
		ok, err := userMgr.IsEmailAvailable(v.User)
		if err != nil {
			return err
		}
		if !ok {
			errors["email"] = "Email already taken"
		}

	}

	if v.User.Password == "" {
		errors["password"] = "Password is missing"
	}

	return nil

}
