package v1

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/infrastracture/persistence"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	repo   *persistence.Repositories
	query  *handler.Server
}

func NewHandler(logger *zap.Logger, repositories *persistence.Repositories, query *handler.Server) *Handler {
	return &Handler{
		logger: logger,
		repo:   repositories,
		query:  query,
	}
}
