package graph

import (
	"context"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/config"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/graph/generated"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/graph/model"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/infrastracture/loader"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/infrastracture/persistence"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/io"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"log"
	"os"
	"strconv"
	"testing"
)

func TestQueryResolver_Users(t *testing.T) {
	t.Run("get all users", func(t *testing.T) {
		resolver := getQueryResolver()
		ctx := context.Background()
		users, err := resolver.Users(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, users)
		assert.Len(t, users, 2)

		want := []*model.User{
			{ID: "2", Name: "Fuga"},
			{ID: "1", Name: "Hoge"},
		}
		if diff := cmp.Diff(want, users); len(diff) != 0 {
			t.Errorf("Users mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestQueryResolver_Messages(t *testing.T) {
	resolver := getQueryResolver()
	ctx := context.Background()
	t.Run("get user=1", func(t *testing.T) {
		messages, err := resolver.Messages(ctx, "1")
		assert.NoError(t, err)
		assert.NotEmpty(t, messages)
		assert.Len(t, messages, 2)

		want := []*model.Message{
			{ID: "2", UserID: "1", Message: "test message id 2"},
			{ID: "1", UserID: "1", Message: "test message id 1"},
		}
		if diff := cmp.Diff(want, messages); len(diff) != 0 {
			t.Errorf("Messages mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("get user=2", func(t *testing.T) {
		messages, err := resolver.Messages(ctx, "2")
		assert.NoError(t, err)
		assert.NotEmpty(t, messages)
		assert.Len(t, messages, 2)

		want := []*model.Message{
			{ID: "4", UserID: "2", Message: "test message id 4"},
			{ID: "3", UserID: "2", Message: "test message id 3"},
		}
		if diff := cmp.Diff(want, messages); len(diff) != 0 {
			t.Errorf("Messages mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("get messages not exist user", func(t *testing.T) {
		messages, err := resolver.Messages(ctx, "9999")
		assert.NoError(t, err)
		assert.NotNil(t, messages)
		assert.Empty(t, messages)
		assert.Len(t, messages, 0)
	})
}

func TestMessageResolver_User(t *testing.T) {
	resolver := getResolver()
	ctx := context.Background()
	msgResolver := resolver.Message()

	t.Run("get user=1", func(t *testing.T) {
		message := model.Message{
			UserID: "1",
		}
		user, err := msgResolver.User(ctx, &message)
		assert.NoError(t, err)
		assert.NotEmpty(t, user)
		assert.Equal(t, "1", user.ID)
		assert.Equal(t, "Hoge", user.Name)
	})
	t.Run("get user=2", func(t *testing.T) {
		message := model.Message{
			UserID: "2",
		}
		user, err := msgResolver.User(ctx, &message)
		assert.NoError(t, err)
		assert.NotEmpty(t, user)
		assert.Equal(t, "2", user.ID)
		assert.Equal(t, "Fuga", user.Name)
	})
	t.Run("get not exist user", func(t *testing.T) {
		message := model.Message{
			UserID: "9999",
		}
		user, err := msgResolver.User(ctx, &message)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestMutationResolver_CreateMessage(t *testing.T) {
	resolver := getMutationResolver()
	ctx := context.Background()

	t.Run("create message user=1", func(t *testing.T) {
		input := model.NewMessage{
			UserID:  "1",
			Message: "new message 1",
		}
		message, err := resolver.CreateMessage(ctx, input)
		assert.NoError(t, err)
		assert.NotEmpty(t, message)
		id, _ := strconv.Atoi(message.ID)
		assert.Greater(t, id, 1)
		assert.Equal(t, input.UserID, message.UserID)
		assert.Equal(t, input.Message, message.Message)

		t.Cleanup(func() {
			if message != nil && message.ID != "0" {
				db := getDatabase()
				st, er := db.Prepare("DELETE FROM message WHERE id = ?")
				if er == nil {
					_, er = st.Exec(message.ID)
					if er != nil {
						panic(er.Error())
					}
				}
				defer st.Close()
			}
		})
	})

	t.Run("create message user=2", func(t *testing.T) {
		input := model.NewMessage{
			UserID:  "2",
			Message: "new message 2",
		}
		message, err := resolver.CreateMessage(ctx, input)
		assert.NoError(t, err)
		assert.NotEmpty(t, message)
		id, _ := strconv.Atoi(message.ID)
		assert.Greater(t, id, 1)
		assert.Equal(t, input.UserID, message.UserID)
		assert.Equal(t, input.Message, message.Message)

		t.Cleanup(func() {
			if message != nil && message.ID != "0" {
				db := getDatabase()
				st, er := db.Prepare("DELETE FROM message WHERE id = ?")
				if er == nil {
					_, er = st.Exec(message.ID)
					if er != nil {
						panic(er.Error())
					}
				}
				defer st.Close()
			}
		})
	})

	t.Run("create not exist user", func(t *testing.T) {
		input := model.NewMessage{
			UserID:  "9999",
			Message: "new message 99999",
		}
		message, err := resolver.CreateMessage(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, message)
	})
	t.Run("create wrong user_id", func(t *testing.T) {
		input := model.NewMessage{
			UserID:  "AAA",
			Message: "new message AAA",
		}
		message, err := resolver.CreateMessage(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, message)
	})
}

func getQueryResolver() generated.QueryResolver {
	return getResolver().Query()
}

func getMutationResolver() generated.MutationResolver {
	return getResolver().Mutation()
}

func getResolver() *Resolver {
	mysqlDatabase := getDatabase()
	repositories, err := persistence.NewRepositories(mysqlDatabase)
	if err != nil {
		log.Println("failed to new repositories", zap.Error(err))
		os.Exit(1)
	}

	loaders := loader.NewLoaders(repositories)
	resolver := &Resolver{
		Repo:    repositories,
		Loaders: loaders,
	}
	return resolver
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
