package validation

import (
	"github.com/danjac/photoshare/api/models"
)

func NewPhotoValidator(photo *models.Photo) *PhotoValidator {
	return &PhotoValidator{photo}
}

type PhotoValidator struct {
	Photo *models.Photo
}

func (v *PhotoValidator) Validate() (*ValidationResult, error) {
	result := NewValidationResult()
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
	return result, nil
}
