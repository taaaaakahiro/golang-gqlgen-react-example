package persistence

import (
	"context"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/config"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/io"
	"log"
	"os"
	"testing"
)

var messageRepo *MessageRepo
var userRepo *UserRepo

func TestMain(m *testing.M) {
	db := getDatabase()
	messageRepo = NewMessageRepository(db)
	userRepo = NewUserRepository(db)
	res := m.Run()
	os.Exit(res)
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
