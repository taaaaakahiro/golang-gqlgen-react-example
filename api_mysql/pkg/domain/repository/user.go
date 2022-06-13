package repository

import (
	"context"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/domain/entity"
	"github.com/graph-gophers/dataloader"
)

type IUserRepository interface {
	ListUsers() ([]entity.User, error)
	User(userId int) (entity.User, error)
	GetUsers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result
}
