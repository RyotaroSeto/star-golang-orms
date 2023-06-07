package pkg

import (
	"github.com/pkg/errors"
)

var (
	ErrNoMorePages = errors.New("no more pages to get")
	ErrNoStars     = errors.New("no stars present")
	ErrOtherReason = errors.New("other reason is error")
	ErrRateLimit   = errors.New("this error is rate limit")
	ErrNotFound    = errors.New("this error is not found")
)
