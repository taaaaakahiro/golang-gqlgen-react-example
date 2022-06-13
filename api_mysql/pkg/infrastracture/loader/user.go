package loader

import (
	"context"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/graph/model"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/infrastracture/persistence"
	"github.com/graph-gophers/dataloader"
)

// UserLoader dataloader for user table
type UserLoader struct {
	loader *dataloader.Loader
}

// NewUserLoader create user Dataloader instance
func NewUserLoader(r *persistence.Repositories) *UserLoader {
	return &UserLoader{
		loader: dataloader.NewBatchedLoader(r.User.GetUsers),
	}
}

// GetUser wraps the User dataloader for efficient retrieval by user ID
func (l *Loaders) GetUser(ctx context.Context, userID string) (*model.User, error) {
	// HACK:
	// userIDをキーにして、オンメモリでキャッシュ。既に該当キーがキャッシュにあればキャッシュから、なければDBから取得している。
	thunk := l.UserLoader.loader.Load(ctx, dataloader.StringKey(userID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result.(*model.User), nil
}
