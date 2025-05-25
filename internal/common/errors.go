package common

import (
	"fmt"
	"runtime"
)

// ErrorType represents the type of error
type ErrorType string

const (
	ErrorTypeConfig     ErrorType = "config"     // configuration errors
	ErrorTypeFileSystem ErrorType = "filesystem" // file system errors
	ErrorTypeValidation ErrorType = "validation" // validation errors
	ErrorTypeRepository ErrorType = "repository" // repository errors
	ErrorTypeService    ErrorType = "service"    // service errors
	ErrorTypeUI         ErrorType = "ui"         // UI errors
)

// AppError represents an application error with context
type AppError struct {
	Type    ErrorType
	Message string
	Cause   error
	File    string
	Line    int
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// New creates a new application error
func New(errType ErrorType, message string) *AppError {
	_, file, line, _ := runtime.Caller(1)
	return &AppError{
		Type:    errType,
		Message: message,
		File:    file,
		Line:    line,
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, errType ErrorType, message string) *AppError {
	_, file, line, _ := runtime.Caller(1)
	return &AppError{
		Type:    errType,
		Message: message,
		Cause:   err,
		File:    file,
		Line:    line,
	}
}
