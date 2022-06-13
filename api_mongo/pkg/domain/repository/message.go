package repository

import (
	"context"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mongo/pkg/domain/entity"
)

type IMessageRepository interface {
	ListMessages(ctx context.Context, userId int) ([]entity.Message, error)
	CreateMessage(ctx context.Context, message *entity.Message) error
}
