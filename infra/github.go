package infra

import (
	"context"
	"star-golang-orms/domain/repository"
)

// //go:generate go run go.uber.org/mock/mockgen -source=github.go -package=resourcemanage -destination=github_mock.go
// type GitHubQuerier interface {
// 	// getStargazersPage(ctx context.Context, repo string, page int) ([]Stargazer, error)
// }

type GitHubRepository struct {
}

var _ repository.GitHub = &GitHubRepository{}

func NewGitHubRepository(ctx context.Context) repository.GitHub {
	return &GitHubRepository{}
}

func (r *GitHubRepository) GetStar(ctx context.Context) (int, error) {
	return 0, nil
}

func (r *GitHubRepository) GetRateLimit(ctx context.Context) (int, error) {
	return 0, nil
}

// func GetRateLimit(token string) error {
// 	url := baseURL + rateLimit
// 	client := NewHttpClient(url, http.MethodGet, token)
// 	res, err := client.SendRequest()
// 	if err != nil {
// 		return err
// 	}

// 	bts, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return err
// 	}
// 	defer res.Body.Close()

// 	switch res.StatusCode {
// 	case http.StatusOK:
// 		var r map[string]interface{}
// 		if err := json.Unmarshal(bts, &r); err != nil {
// 			return err
// 		}
// 		fmt.Println(r)
// 		return nil
// 	case http.StatusNotModified:
// 		return ErrRateLimit
// 	case http.StatusNotFound:
// 		return ErrNotFound
// 	default:
// 		return ErrOtherReason
// 	}
// }
