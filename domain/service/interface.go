package service

import (
	"context"
)

//go:generate go run go.uber.org/mock/mockgen -source=interface.go -package=service -destination=interface_mock.go
type Fetcher interface {
	Start(ctx context.Context) error
}
