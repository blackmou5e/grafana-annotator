package errors

import "fmt"

type ErrorType string

const (
	ErrConfiguration ErrorType = "ConfigurationError"
	ErrGrafanaAPI    ErrorType = "GrafanaAPIError"
	ErrValidation    ErrorType = "ValidationError"
	ErrInternal      ErrorType = "InternalError"
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func NewAppError(errType ErrorType, message string, err error) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}
