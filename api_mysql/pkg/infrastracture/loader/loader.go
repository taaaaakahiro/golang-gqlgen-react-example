package loader

import (
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/infrastracture/persistence"
)

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	UserLoader *UserLoader
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(r *persistence.Repositories) *Loaders {
	// define the data loader
	loaders := &Loaders{
		UserLoader: NewUserLoader(r),
	}
	return loaders
}
