package errors

import (
	"fmt"
	"net/http"
)

const (
	BadRequest          = http.StatusBadRequest
	UnprocessableEntity = http.StatusUnprocessableEntity
	Forbidden           = http.StatusForbidden
	NotFound            = http.StatusNotFound
	RequestTimeout      = http.StatusRequestTimeout
	Conflict            = http.StatusConflict
	Unauthorized        = http.StatusUnauthorized
	InternalServerError = http.StatusInternalServerError
)

var (
	ErrNoMorePages = New(NotFound, "no more pages to get")
	ErrNoStars     = New(InternalServerError, "no stars present")
	ErrOtherReason = New(Unauthorized, "other reason is error")
	ErrRateLimit   = New(InternalServerError, "this error is rate limit")
	ErrNotFound    = New(NotFound, "this error is not found")
)

type CustomError interface {
	error
	Code() int
}

type customError struct {
	code int
	err  error
	msg  string
}

func (e *customError) Error() string {
	if e == nil {
		return ""
	}

	return e.msg
}

func New(code int, msg string) CustomError {
	return &customError{code: code, err: nil, msg: msg}
}

func Newf(code int, format string, args ...interface{}) CustomError {
	return &customError{code: code, err: nil, msg: fmt.Sprintf(format, args...)}
}

func (e *customError) Code() int {
	if e == nil {
		return http.StatusOK
	}

	return e.code
}
