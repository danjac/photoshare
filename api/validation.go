package api

import (
	"regexp"
)

var emailRegex = regexp.MustCompile(".+@.+\\..+")

type validationFailure struct {
	Errors map[string]string `json:"errors"`
}

func (f validationFailure) Error() string {
	return "Validation failure"
}

func validate(v validator) error {
	errors := make(map[string]string)
	if err := v.validate(errors); err != nil {
		return err
	}
	if len(errors) > 0 {
		return validationFailure{errors}
	}
	return nil
}

type validator interface {
	validate(map[string]string) error
}

func newPhotoValidator(photo *photo) *photoValidator {
	return &photoValidator{photo}
}

type photoValidator struct {
	photo *photo
}

func (v *photoValidator) validate(errors map[string]string) error {
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

func newUserValidator(user *user, ds userDataStore) *userValidator {
	return &userValidator{user, ds}
}

type userValidator struct {
	user   *user
	userDS userDataStore
}

func (v *userValidator) validate(errors map[string]string) error {

	if v.user.Name == "" {
		errors["name"] = "Name is missing"
	} else {
		ok, err := v.userDS.isNameAvailable(v.user)
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
		ok, err := v.userDS.isEmailAvailable(v.user)
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
