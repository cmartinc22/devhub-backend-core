package models

import (
	"fmt"
)

type ErrorCode string

func (err ErrorCode) String() string {
	return string(err)
}

const (
	PreconditionFailed   ErrorCode = "PRECONDITION_FAILED"
	InvalidParameters    ErrorCode = "INVALID_PARAMETERS"
	UnsupportedMediaType ErrorCode = "UNSOPPORTED_MEDIA_TYPE"
	NotFound             ErrorCode = "NOT_FOUND"
	InternalError        ErrorCode = "INTERNAL_SERVER_ERROR"
	IncorrectHeaders     ErrorCode = "INCORRECT_HEADERS"
	InvalidCursor        ErrorCode = "INVALID_CURSOR"
	Forbidden            ErrorCode = "FORBIDDEN"
)

const (
	UnsupportedVersion      ErrorCode = "UNSUPPORTED_API_VERSION"
	LesserVersionRequested  ErrorCode = "LESSER_VERSION_REQUESTED"
	GreaterVersionRequested ErrorCode = "GREATER_VERSION_REQUESTED"
	InferenceError          ErrorCode = "INFERENCE_ERROR"
	SpecValidateError       ErrorCode = "SPEC_VALIDATION_ERROR"
)

type CustomError struct {
	Code     ErrorCode
	Messages []string
}

type EntityError struct {
	Path    string    `json:"path" yaml:"path"`
	Code    ErrorCode `json:"code" yaml:"code"`
	Message string    `json:"message" yaml:"message"`
}

func (err CustomError) String() string {
	return fmt.Sprintf("ERROR[%s] %v", err.Code, err.Messages)
}

func (err *CustomError) Error() string {
	return err.String()
}
