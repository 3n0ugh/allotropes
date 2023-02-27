package errors

import (
	"net/http"

	"github.com/stretchr/testify/assert"
)

func NewBadRequestError(message, devMessage string) *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,
		Title:      "bad request",
		Message:    message,
		devMessage: devMessage,
	}
}

func NewInternalServerError(devMessage string) *Error {
	return &Error{
		StatusCode: http.StatusInternalServerError,
		Title:      "internal server error",
		Message:    "something goes wrong",
		devMessage: devMessage,
	}
}

func NewNotFoundError(message, devMessage string) *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		Title:      "not found",
		Message:    message,
		devMessage: devMessage,
	}
}

func NewForbiddenError(message, devMessage string) *Error {
	return &Error{
		StatusCode: http.StatusForbidden,
		Title:      "forbidden",
		Message:    message,
		devMessage: devMessage,
	}
}

func NewUnAuthorizedError(message, devMessage string) *Error {
	return &Error{
		StatusCode: http.StatusUnauthorized,
		Title:      "unauthorized",
		Message:    message,
		devMessage: devMessage,
	}
}

func CompareError(e1, e2 *Error) bool {
	return assert.ObjectsAreEqual(e1, e2)
}
