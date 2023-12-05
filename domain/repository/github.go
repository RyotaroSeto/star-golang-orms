package repository

import "context"

//go:generate go run go.uber.org/mock/mockgen -source=github.go -package=repository -destination=github_mock.go
type GitHub interface {
	GetStar(ctx context.Context) (int, error)
	GetRateLimit(ctx context.Context) (int, error)
}
