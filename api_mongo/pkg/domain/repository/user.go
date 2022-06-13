package repository

import (
	"context"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mongo/pkg/domain/entity"
)

type IUserRepository interface {
	ListUsers(ctx context.Context) ([]entity.User, error)
	User(ctx context.Context, userId int) (entity.User, error)
}
