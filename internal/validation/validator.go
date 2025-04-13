package validation

import (
	"strings"

	"github.com/blackmou5e/grafana-annotator/pkg/errors"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateAnnotationInput(tags []string, message string) error {
	if len(tags) == 0 {
		return errors.NewAppError(errors.ErrValidation, "At least one tag is required", nil)
	}

	for _, tag := range tags {
		if strings.TrimSpace(tag) == "" {
			return errors.NewAppError(errors.ErrValidation, "Empty tags are not allowed", nil)
		}
	}

	if strings.TrimSpace(message) == "" {
		return errors.NewAppError(errors.ErrValidation, "Annotation message cannot be empty", nil)
	}

	return nil
}
