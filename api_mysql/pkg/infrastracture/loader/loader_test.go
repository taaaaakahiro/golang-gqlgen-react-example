package loader

import (
	"context"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/config"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/infrastracture/persistence"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/io"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"log"
	"os"
	"testing"
)

var testLoaders *Loaders

func TestMain(m *testing.M) {
	db := getDatabase()
	repositories, err := persistence.NewRepositories(db)
	if err != nil {
		log.Println("failed to new repositories", zap.Error(err))
		os.Exit(1)
	}
	testLoaders = NewLoaders(repositories)
	res := m.Run()
	os.Exit(res)
}

func TestNewLoaders(t *testing.T) {
	assert.NotNil(t, testLoaders.UserLoader)
}

func getDatabase() *io.SQLDatabase {
	// prepare db
	cfg, err := config.LoadConfig(context.Background())
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	sqlSetting := &config.SQLDBSettings{
		SqlDsn:              cfg.DB.DSN,
		SqlMaxOpenConns:     cfg.DB.MaxOpenConns,
		SqlMaxIdleConns:     cfg.DB.MaxIdleConns,
		SqlConnsMaxLifetime: cfg.DB.ConnsMaxLifetime,
	}

	db, err := io.NewDatabase(sqlSetting)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	return db
}
