package repository

import (
	"context"
	"star-golang-orms/domain/model"
)

//go:generate go run go.uber.org/mock/mockgen -source=github.go -package=repository -destination=github_mock.go
type GitHub interface {
	GetRepository(ctx context.Context, rn model.RepositoryName) (*model.Repository, error)
	GetStarPage(ctx context.Context, repo *model.Repository, page int) (*[]model.Stargazer, error)
}
