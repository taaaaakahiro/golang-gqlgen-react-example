package persistence

import (
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/domain/repository"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/io"
)

type Repositories struct {
	User    repository.IUserRepository
	Message repository.IMessageRepository
}

func NewRepositories(db *io.SQLDatabase) (*Repositories, error) {
	return &Repositories{
		User:    NewUserRepository(db),
		Message: NewMessageRepository(db),
	}, nil
}
