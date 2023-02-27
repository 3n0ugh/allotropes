package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

type Error struct {
	StatusCode int    `json:"statusCode"`
	Title      string `json:"title"`
	Message    string `json:"message"`
	devMessage string `json:"-"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s", e.devMessage)
}

func New(str string) error {
	return errors.New(str)
}

func Wrap(err error, str string) error {
	return errors.Wrap(err, str)
}
