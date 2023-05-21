package pkg

import (
	"github.com/pkg/errors"
)

var (
	ErrNoMorePages  = errors.New("no more pages to get")
	ErrTooManyStars = errors.New("repo has too many stargazers, github won't allow us to list all stars")
)

func ValidateStatusCode(statusCode int) error {
	switch statusCode {
	case 304:
		return errors.Errorf("failed to not modified: %s", statusCode)
	case 401:
		return errors.Errorf("failed to requires authentication: %s", statusCode)
	case 403:
		return errors.Errorf("failed to forbidden: %s", statusCode)
	case 422:
		return errors.Errorf("failed to endpoint has been spammed: %s", statusCode)
	}
	return nil
}
