package repository

import (
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/domain/entity"
)

type IMessageRepository interface {
	ListMessages(userId int) ([]entity.Message, error)
	Messages() ([]entity.Message, error)
	CreateMessage(message *entity.Message) error
}
