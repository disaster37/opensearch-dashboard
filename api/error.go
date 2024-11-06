package api

import (
	"fmt"
)

// APIError is the error object
type APIError struct {
	Code    int
	Message string
}

// Error return error message
func (e APIError) Error() string {
	return e.Message
}

// NewAPIErrorf create new API error with code and message
func NewAPIErrorf(code int, message string, params ...interface{}) APIError {

	return APIError{
		Code:    code,
		Message: fmt.Sprintf(message, params...),
	}

}

// NewAPIErrorf create new API error with code and message
func NewAPIError(code int, message string) APIError {

	return APIError{
		Code:    code,
		Message: message,
	}

}
