package graph

import (
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/infrastracture/loader"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/infrastracture/persistence"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repo    *persistence.Repositories
	Loaders *loader.Loaders
}
