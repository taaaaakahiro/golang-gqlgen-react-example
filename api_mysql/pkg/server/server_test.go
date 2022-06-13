package server

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/99designs/gqlgen/graphql"
	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/config"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/domain/entity"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/graph"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/graph/generated"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/graph/model"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/handler"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/infrastracture/loader"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/infrastracture/persistence"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/io"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func TestServer(t *testing.T) {
	// create logger
	logger, err := zap.NewProduction()
	if err != nil {
		t.Errorf("failed to setup loggerr: %s\n", err)
	}
	defer logger.Sync()

	// load config
	ctx := context.Background()
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		t.Errorf("failed to load config: %s\n", err)
	}

	// init mysql
	logger.Info("connect to mysql ", zap.String("DSN", cfg.DB.DSN))
	sqlSetting := &config.SQLDBSettings{
		SqlDsn:              cfg.DB.DSN,
		SqlMaxOpenConns:     cfg.DB.MaxOpenConns,
		SqlMaxIdleConns:     cfg.DB.MaxIdleConns,
		SqlConnsMaxLifetime: cfg.DB.ConnsMaxLifetime,
	}

	mysqlDatabase, err := io.NewDatabase(sqlSetting)
	if err != nil {
		t.Errorf("failed to create mysql db repository: %s\n", err)
	}

	repositories, err := persistence.NewRepositories(mysqlDatabase)
	assert.NoError(t, err)

	// start server

	loaders := loader.NewLoaders(repositories)
	// init to start http server
	// init gql server
	query := gqlhandler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{
			Resolvers: &graph.Resolver{
				Repo:    repositories,
				Loaders: loaders,
			},
		}))
	registry := handler.NewHandler(logger, repositories, query, "v1.0-test")
	s := NewServer(registry, &Config{Log: logger})
	testServer := httptest.NewServer(s.Mux)
	defer testServer.Close()

	// start API test
	t.Run("check /healthz", func(t *testing.T) {
		res, err := http.Get(testServer.URL + "/healthz")
		assert.NoError(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)
	})

	t.Run("check /version", func(t *testing.T) {
		res, err := http.Get(testServer.URL + "/version")
		assert.NoError(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)

		// read body
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		assert.NoError(t, err)
		assert.NotEmpty(t, body)

		var data interface{}
		err = json.Unmarshal(body, &data)
		assert.NoError(t, err)
		ver := data.(map[string]interface{})["version"].(string)
		assert.Equal(t, ver, "v1.0-test")
	})

	t.Run("check /user/list", func(t *testing.T) {
		res, err := http.Get(testServer.URL + "/user/list")
		assert.NoError(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)

		// read body
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		assert.NoError(t, err)
		assert.NotEmpty(t, body)

		var users []entity.User
		err = json.Unmarshal(body, &users)
		assert.NoError(t, err)
		assert.Greater(t, len(users), 0)

		assert.Contains(t, users, entity.User{
			Id:   1,
			Name: "Hoge",
		})
		assert.Contains(t, users, entity.User{
			Id:   2,
			Name: "Fuga",
		})
	})

	t.Run("check /message/1", func(t *testing.T) {
		res, err := http.Get(testServer.URL + "/message/list/1")
		assert.NoError(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)

		// read body
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		assert.NoError(t, err)
		assert.NotEmpty(t, body)

		var messages []entity.Message
		err = json.Unmarshal(body, &messages)
		assert.NoError(t, err)
		assert.Greater(t, len(messages), 0)

		assert.Contains(t, messages, entity.Message{
			Id:      1,
			UserId:  1,
			Message: "test message id 1",
		})
		assert.Contains(t, messages, entity.Message{
			Id:      2,
			UserId:  1,
			Message: "test message id 2",
		})
	})

	t.Run("check /message/2", func(t *testing.T) {
		res, err := http.Get(testServer.URL + "/message/list/2")
		assert.NoError(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)

		// read body
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		assert.NoError(t, err)
		assert.NotEmpty(t, body)

		var messages []entity.Message
		err = json.Unmarshal(body, &messages)
		assert.NoError(t, err)
		assert.Greater(t, len(messages), 0)

		assert.Contains(t, messages, entity.Message{
			Id:      3,
			UserId:  2,
			Message: "test message id 3",
		})
		assert.Contains(t, messages, entity.Message{
			Id:      4,
			UserId:  2,
			Message: "test message id 4",
		})
	})

	// GraphQL endpoints
	t.Run("check /gql", func(t *testing.T) {
		res, err := http.Get(testServer.URL + "/gql")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
	t.Run("check /query", func(t *testing.T) {
		res, err := http.Get(testServer.URL + "/query")
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusOK, res.StatusCode)

		// read body
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		assert.NoError(t, err)
		assert.NotEmpty(t, body)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		assert.NoError(t, err)

		assert.NotEmpty(t, result["errors"])
		assert.Nil(t, result["data"])
	})
	t.Run("check graphQL query", func(t *testing.T) {

		t.Run("query users", func(t *testing.T) {
			type resp struct {
				Users []model.User `json:"users"`
			}

			t.Run("get id, name", func(t *testing.T) {
				graphReq := graphql.RawParams{
					Query: `
query getUsers {
 users {
   id
   name
 }
}
			`,
				}
				req, _ := json.Marshal(graphReq)

				res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
				assert.NoError(t, err)
				// read body
				body, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				assert.NoError(t, err)
				assert.NotEmpty(t, body)

				var gqlRes graphql.Response
				err = json.Unmarshal(body, &gqlRes)
				assert.NoError(t, err)
				assert.Nil(t, gqlRes.Errors)
				assert.NotEmpty(t, gqlRes.Data)

				var data resp
				err = json.Unmarshal(gqlRes.Data, &data)
				assert.NoError(t, err)
				assert.NotNil(t, data.Users)

				want := []model.User{
					{ID: "2", Name: "Fuga"},
					{ID: "1", Name: "Hoge"},
				}
				if diff := cmp.Diff(want, data.Users); len(diff) != 0 {
					t.Errorf("Users mismatch (-want +got):\n%s", diff)
				}
			})

			t.Run("get id", func(t *testing.T) {
				graphReq := graphql.RawParams{
					Query: `
query getUsers {
 users {
   id
 }
}
			`,
				}
				req, _ := json.Marshal(graphReq)

				res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
				assert.NoError(t, err)
				// read body
				body, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				assert.NoError(t, err)
				assert.NotEmpty(t, body)

				var gqlRes graphql.Response
				err = json.Unmarshal(body, &gqlRes)
				assert.NoError(t, err)
				assert.Nil(t, gqlRes.Errors)
				assert.NotEmpty(t, gqlRes.Data)

				var data resp
				err = json.Unmarshal(gqlRes.Data, &data)
				assert.NoError(t, err)
				assert.NotNil(t, data.Users)

				want := []model.User{
					{ID: "2"},
					{ID: "1"},
				}
				if diff := cmp.Diff(want, data.Users); len(diff) != 0 {
					t.Errorf("Users mismatch (-want +got):\n%s", diff)
				}
			})

			t.Run("get name", func(t *testing.T) {
				graphReq := graphql.RawParams{
					Query: `
query getUsers {
 users {
   name
 }
}
			`,
				}
				req, _ := json.Marshal(graphReq)

				res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
				assert.NoError(t, err)
				// read body
				body, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				assert.NoError(t, err)
				assert.NotEmpty(t, body)

				var gqlRes graphql.Response
				err = json.Unmarshal(body, &gqlRes)
				assert.NoError(t, err)
				assert.Nil(t, gqlRes.Errors)
				assert.NotEmpty(t, gqlRes.Data)

				var data resp
				err = json.Unmarshal(gqlRes.Data, &data)
				assert.NoError(t, err)
				assert.NotNil(t, data.Users)

				want := []model.User{
					{Name: "Fuga"},
					{Name: "Hoge"},
				}
				if diff := cmp.Diff(want, data.Users); len(diff) != 0 {
					t.Errorf("Users mismatch (-want +got):\n%s", diff)
				}
			})
		})

		t.Run("query messages", func(t *testing.T) {
			type respMessage struct {
				ID      string     `json:"id"`
				Message string     `json:"message"`
				User    model.User `json:"user"`
			}

			type resp struct {
				Messages []respMessage `json:"messages"`
			}

			t.Run("get messages user=1 with user id,name", func(t *testing.T) {
				graphReq := graphql.RawParams{
					Query: `
query getMessages($userID: ID!) {
  messages(userID: $userID) {
    id
    message
    user {
      id
      name
    }
  }
}
`,
					Variables: map[string]interface{}{
						"userID": 1,
					},
				}
				req, _ := json.Marshal(graphReq)

				res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
				assert.NoError(t, err)
				// read body
				body, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				assert.NoError(t, err)
				assert.NotEmpty(t, body)

				var gqlRes graphql.Response
				err = json.Unmarshal(body, &gqlRes)
				assert.NoError(t, err)
				assert.Nil(t, gqlRes.Errors)
				assert.NotEmpty(t, gqlRes.Data)

				var data resp
				err = json.Unmarshal(gqlRes.Data, &data)
				assert.NoError(t, err)
				assert.NotNil(t, data.Messages)

				want := []respMessage{
					{ID: "2", Message: "test message id 2", User: model.User{ID: "1", Name: "Hoge"}},
					{ID: "1", Message: "test message id 1", User: model.User{ID: "1", Name: "Hoge"}},
				}
				if diff := cmp.Diff(want, data.Messages); len(diff) != 0 {
					t.Errorf("Messages mismatch (-want +got):\n%s", diff)
				}
			})

			t.Run("get messages user=2 without id", func(t *testing.T) {
				graphReq := graphql.RawParams{
					Query: `
query getMessages($userID: ID!) {
  messages(userID: $userID) {
    message
    user {
      name
    }
  }
}`,
					Variables: map[string]interface{}{
						"userID": 2,
					},
				}
				req, _ := json.Marshal(graphReq)

				res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
				assert.NoError(t, err)
				// read body
				body, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				assert.NoError(t, err)
				assert.NotEmpty(t, body)

				var gqlRes graphql.Response
				err = json.Unmarshal(body, &gqlRes)
				assert.NoError(t, err)
				assert.Nil(t, gqlRes.Errors)
				assert.NotEmpty(t, gqlRes.Data)

				var data resp
				err = json.Unmarshal(gqlRes.Data, &data)
				assert.NoError(t, err)
				assert.NotNil(t, data.Messages)

				want := []respMessage{
					{Message: "test message id 4", User: model.User{Name: "Fuga"}},
					{Message: "test message id 3", User: model.User{Name: "Fuga"}},
				}
				if diff := cmp.Diff(want, data.Messages); len(diff) != 0 {
					t.Errorf("Messages mismatch (-want +got):\n%s", diff)
				}
			})

			t.Run("get messages user=2 without users", func(t *testing.T) {
				graphReq := graphql.RawParams{
					Query: `
query getMessages($userID: ID!) {
  messages(userID: $userID) {
    message
  }
}`,
					Variables: map[string]interface{}{
						"userID": 2,
					},
				}
				req, _ := json.Marshal(graphReq)

				res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
				assert.NoError(t, err)
				// read body
				body, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				assert.NoError(t, err)
				assert.NotEmpty(t, body)

				var gqlRes graphql.Response
				err = json.Unmarshal(body, &gqlRes)
				assert.NoError(t, err)
				assert.Nil(t, gqlRes.Errors)
				assert.NotEmpty(t, gqlRes.Data)

				var data resp
				err = json.Unmarshal(gqlRes.Data, &data)
				assert.NoError(t, err)
				assert.NotNil(t, data.Messages)

				want := []respMessage{
					{Message: "test message id 4"},
					{Message: "test message id 3"},
				}
				if diff := cmp.Diff(want, data.Messages); len(diff) != 0 {
					t.Errorf("Messages mismatch (-want +got):\n%s", diff)
				}
			})

			t.Run("get messages not exist user", func(t *testing.T) {
				graphReq := graphql.RawParams{
					Query: `
query getMessages($userID: ID!) {
  messages(userID: $userID) {
    id
    message
    user {
      id
      name
    }
  }
}`,
					Variables: map[string]interface{}{
						"userID": 9999,
					},
				}
				req, _ := json.Marshal(graphReq)

				res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
				assert.NoError(t, err)
				// read body
				body, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				assert.NoError(t, err)
				assert.NotEmpty(t, body)

				var gqlRes graphql.Response
				err = json.Unmarshal(body, &gqlRes)
				assert.NoError(t, err)
				assert.Nil(t, gqlRes.Errors)
				assert.NotEmpty(t, gqlRes.Data)

				var data resp
				err = json.Unmarshal(gqlRes.Data, &data)
				assert.NoError(t, err)
				assert.NotNil(t, data.Messages)
				assert.Empty(t, data.Messages)
			})

			t.Run("invalid param", func(t *testing.T) {
				graphReq := graphql.RawParams{
					Query: `query getMessages($userID: ID!) {
				 messages(userID: $userID) {
				   id
				   message
				   user {
				     id
				     name
				   }
				 }
				}`,
					Variables: map[string]interface{}{
						"userID": "AAAAA",
					},
				}
				req, _ := json.Marshal(graphReq)

				res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
				assert.NoError(t, err)
				// read body
				body, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				assert.NoError(t, err)
				assert.NotEmpty(t, body)

				var gqlRes graphql.Response
				err = json.Unmarshal(body, &gqlRes)
				assert.NoError(t, err)
				assert.NotNil(t, gqlRes.Errors)
				assert.NotEmpty(t, gqlRes.Errors)
				assert.NotNil(t, gqlRes.Data)

				var data resp
				err = json.Unmarshal(gqlRes.Data, &data)
				assert.NoError(t, err)
				assert.Nil(t, data.Messages)
			})

			t.Run("lack params", func(t *testing.T) {
				graphReq := graphql.RawParams{
					Query: `
query getMessages($userID: ID!) {
 messages(userID: $userID) {
   id
   message
   user {
	 id
	 name
   }
 }
}
`,
				}
				req, _ := json.Marshal(graphReq)

				res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
				assert.NoError(t, err)
				// read body
				body, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				assert.NoError(t, err)
				assert.NotEmpty(t, body)

				var gqlRes graphql.Response
				err = json.Unmarshal(body, &gqlRes)
				assert.NoError(t, err)
				assert.NotNil(t, gqlRes.Errors)
				assert.NotEmpty(t, gqlRes.Errors)
				assert.NotNil(t, gqlRes.Data)

				var data resp
				err = json.Unmarshal(gqlRes.Data, &data)
				assert.NoError(t, err)
				assert.Nil(t, data.Messages)
			})
		})

		t.Run("query all messages", func(t *testing.T) {
			type respMessage struct {
				ID      string     `json:"id"`
				Message string     `json:"message"`
				User    model.User `json:"user"`
			}

			type resp struct {
				AllMessages []respMessage `json:"allMessages"`
			}

			t.Run("get all messages with user id,name", func(t *testing.T) {
				graphReq := graphql.RawParams{
					Query: `
query getAllMessages {
  allMessages {
    id
    message
    user {
      id
      name
    }
  }
}
`,
				}
				req, _ := json.Marshal(graphReq)

				res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
				assert.NoError(t, err)
				// read body
				body, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				assert.NoError(t, err)
				assert.NotEmpty(t, body)

				var gqlRes graphql.Response
				err = json.Unmarshal(body, &gqlRes)
				assert.NoError(t, err)
				assert.Nil(t, gqlRes.Errors)
				assert.NotEmpty(t, gqlRes.Data)

				var data resp
				err = json.Unmarshal(gqlRes.Data, &data)
				assert.NoError(t, err)
				assert.NotNil(t, data.AllMessages)

				want := []respMessage{
					{ID: "4", Message: "test message id 4", User: model.User{ID: "2", Name: "Fuga"}},
					{ID: "3", Message: "test message id 3", User: model.User{ID: "2", Name: "Fuga"}},
					{ID: "2", Message: "test message id 2", User: model.User{ID: "1", Name: "Hoge"}},
					{ID: "1", Message: "test message id 1", User: model.User{ID: "1", Name: "Hoge"}},
				}
				if diff := cmp.Diff(want, data.AllMessages); len(diff) != 0 {
					t.Errorf("Messages mismatch (-want +got):\n%s", diff)
				}
			})

			t.Run("get all messages without id", func(t *testing.T) {
				graphReq := graphql.RawParams{
					Query: `
query getAllMessages {
  allMessages {
    message
    user {
      name
    }
  }
}
`,
				}
				req, _ := json.Marshal(graphReq)

				res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
				assert.NoError(t, err)
				// read body
				body, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				assert.NoError(t, err)
				assert.NotEmpty(t, body)

				var gqlRes graphql.Response
				err = json.Unmarshal(body, &gqlRes)
				assert.NoError(t, err)
				assert.Nil(t, gqlRes.Errors)
				assert.NotEmpty(t, gqlRes.Data)

				var data resp
				err = json.Unmarshal(gqlRes.Data, &data)
				assert.NoError(t, err)
				assert.NotNil(t, data.AllMessages)

				want := []respMessage{
					{Message: "test message id 4", User: model.User{Name: "Fuga"}},
					{Message: "test message id 3", User: model.User{Name: "Fuga"}},
					{Message: "test message id 2", User: model.User{Name: "Hoge"}},
					{Message: "test message id 1", User: model.User{Name: "Hoge"}},
				}
				if diff := cmp.Diff(want, data.AllMessages); len(diff) != 0 {
					t.Errorf("Messages mismatch (-want +got):\n%s", diff)
				}
			})

			t.Run("get all messages without users", func(t *testing.T) {
				graphReq := graphql.RawParams{
					Query: `
query getAllMessages {
  allMessages {
    message
  }
}
`,
				}
				req, _ := json.Marshal(graphReq)

				res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
				assert.NoError(t, err)
				// read body
				body, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				assert.NoError(t, err)
				assert.NotEmpty(t, body)

				var gqlRes graphql.Response
				err = json.Unmarshal(body, &gqlRes)
				assert.NoError(t, err)
				assert.Nil(t, gqlRes.Errors)
				assert.NotEmpty(t, gqlRes.Data)

				var data resp
				err = json.Unmarshal(gqlRes.Data, &data)
				assert.NoError(t, err)
				assert.NotNil(t, data.AllMessages)

				want := []respMessage{
					{Message: "test message id 4"},
					{Message: "test message id 3"},
					{Message: "test message id 2"},
					{Message: "test message id 1"},
				}
				if diff := cmp.Diff(want, data.AllMessages); len(diff) != 0 {
					t.Errorf("Messages mismatch (-want +got):\n%s", diff)
				}
			})
		})

		t.Run("query multiple models", func(t *testing.T) {
			type respMessage struct {
				ID      string     `json:"id"`
				Message string     `json:"message"`
				User    model.User `json:"user"`
			}

			type resp struct {
				Messages    []respMessage `json:"messages"`
				Users       []model.User  `json:"users"`
				AllMessages []respMessage `json:"allMessages"`
			}

			t.Run("get user=1", func(t *testing.T) {
				graphReq := graphql.RawParams{
					Query: `
query getMessagesAndUsers($userID: ID!) {
  messages(userID: $userID) {
    id
    message
    user {
      id
      name
    }
  }
  users {
    id
    name
  }
  allMessages {
    id
    message
    user {
      id
      name
    }
  }
}
`,
					Variables: map[string]interface{}{
						"userID": 1,
					},
				}
				req, _ := json.Marshal(graphReq)

				res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
				assert.NoError(t, err)
				// read body
				body, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				assert.NoError(t, err)
				assert.NotEmpty(t, body)

				var gqlRes graphql.Response
				err = json.Unmarshal(body, &gqlRes)
				assert.NoError(t, err)
				assert.Nil(t, gqlRes.Errors)
				assert.NotEmpty(t, gqlRes.Data)

				var data resp
				err = json.Unmarshal(gqlRes.Data, &data)
				assert.NoError(t, err)
				assert.NotNil(t, data.Messages)

				wantMsg := []respMessage{
					{ID: "2", Message: "test message id 2", User: model.User{ID: "1", Name: "Hoge"}},
					{ID: "1", Message: "test message id 1", User: model.User{ID: "1", Name: "Hoge"}},
				}
				if diff := cmp.Diff(wantMsg, data.Messages); len(diff) != 0 {
					t.Errorf("Messages mismatch (-want +got):\n%s", diff)
				}
				wantUsr := []model.User{
					{ID: "2", Name: "Fuga"},
					{ID: "1", Name: "Hoge"},
				}
				if diff := cmp.Diff(wantUsr, data.Users); len(diff) != 0 {
					t.Errorf("Users mismatch (-want +got):\n%s", diff)
				}

				wantAllMsg := []respMessage{
					{ID: "4", Message: "test message id 4", User: model.User{ID: "2", Name: "Fuga"}},
					{ID: "3", Message: "test message id 3", User: model.User{ID: "2", Name: "Fuga"}},
					{ID: "2", Message: "test message id 2", User: model.User{ID: "1", Name: "Hoge"}},
					{ID: "1", Message: "test message id 1", User: model.User{ID: "1", Name: "Hoge"}},
				}
				if diff := cmp.Diff(wantAllMsg, data.AllMessages); len(diff) != 0 {
					t.Errorf("AllMessages mismatch (-want +got):\n%s", diff)
				}
			})
		})

		t.Run("undefined schema", func(t *testing.T) {
			graphReq := graphql.RawParams{
				Query: `
query getWrong() {
 wrong
}
`,
			}
			req, _ := json.Marshal(graphReq)

			res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
			assert.NoError(t, err)
			// read body
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			assert.NoError(t, err)
			assert.NotEmpty(t, body)

			var gqlRes graphql.Response
			err = json.Unmarshal(body, &gqlRes)
			assert.NoError(t, err)
			assert.NotNil(t, gqlRes.Errors)
			assert.NotEmpty(t, gqlRes.Errors)
			assert.NotNil(t, gqlRes.Data)

			var data interface{}
			err = json.Unmarshal(gqlRes.Data, &data)
			assert.NoError(t, err)
			assert.Nil(t, data)
		})
	})

	t.Run("check graphQL mutation", func(t *testing.T) {
		type respMessage struct {
			ID      string     `json:"id"`
			Message string     `json:"message"`
			User    model.User `json:"user"`
		}

		type resp struct {
			CreateMessage respMessage `json:"createMessage"`
		}

		t.Run("create messages user=1", func(t *testing.T) {
			graphReq := graphql.RawParams{
				Query: `
mutation createMessage($input: NewMessage!) {
  createMessage(input: $input) {
    id
    message
    user {
      id
      name
    }
  }
}
`,
				Variables: map[string]interface{}{
					"input": map[string]interface{}{
						"userID":  1,
						"message": "new message mutation 111",
					},
				},
			}
			req, _ := json.Marshal(graphReq)

			res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
			assert.NoError(t, err)
			// read body
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			assert.NoError(t, err)
			assert.NotEmpty(t, body)

			var gqlRes graphql.Response
			err = json.Unmarshal(body, &gqlRes)
			assert.NoError(t, err)
			assert.Nil(t, gqlRes.Errors)
			assert.NotEmpty(t, gqlRes.Data)

			var data resp
			err = json.Unmarshal(gqlRes.Data, &data)
			assert.NoError(t, err)
			assert.NotNil(t, data.CreateMessage)

			want := respMessage{Message: "new message mutation 111", User: model.User{ID: "1", Name: "Hoge"}}

			if diff := cmp.Diff(want, data.CreateMessage, cmpopts.IgnoreFields(respMessage{}, "ID")); len(diff) != 0 {
				t.Errorf("Messages mismatch (-want +got):\n%s", diff)
			}

			log.Println()
			log.Printf("data.CreateMessage %+v", data.CreateMessage)
			log.Println()

			t.Cleanup(func() {

				//if data.CreateMessage.ID != "" {
				db := getDatabase()
				st, _ := db.Prepare("DELETE FROM message WHERE id = ?")
				id, _ := strconv.Atoi(data.CreateMessage.ID)
				_, er := st.Exec(id)
				if er != nil {
					panic(er.Error())
				}
				st.Close()
				//}
			})
		})

		t.Run("invalid params (no param)", func(t *testing.T) {
			graphReq := graphql.RawParams{
				Query: `
mutation createMessage($input: NewMessage!) {
  createMessage(input: $input) {
    id
    message
    user {
      id
      name
    }
  }
}
`,
			}
			req, _ := json.Marshal(graphReq)

			res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
			assert.NoError(t, err)
			// read body
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			assert.NoError(t, err)
			assert.NotEmpty(t, body)

			var gqlRes graphql.Response
			err = json.Unmarshal(body, &gqlRes)
			assert.NoError(t, err)
			assert.NotEmpty(t, gqlRes.Errors)
			assert.NotNil(t, gqlRes.Data)

			var data interface{}
			err = json.Unmarshal(gqlRes.Data, &data)
			assert.Nil(t, data)
		})

		t.Run("invalid params (blank input)", func(t *testing.T) {
			graphReq := graphql.RawParams{
				Query: `
mutation createMessage($input: NewMessage!) {
 createMessage(input: $input) {
   id
   message
   user {
	 id
	 name
   }
 }
}
		`,
				Variables: map[string]interface{}{
					"input": map[string]interface{}{},
				},
			}
			req, _ := json.Marshal(graphReq)

			res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
			assert.NoError(t, err)
			// read body
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			assert.NoError(t, err)
			assert.NotEmpty(t, body)

			var gqlRes graphql.Response
			err = json.Unmarshal(body, &gqlRes)
			assert.NoError(t, err)
			assert.NotEmpty(t, gqlRes.Errors)
			assert.NotEmpty(t, gqlRes.Data)

			var data interface{}
			err = json.Unmarshal(gqlRes.Data, &data)
			assert.Nil(t, data)
		})

		t.Run("invalid params (wrong format)", func(t *testing.T) {
			graphReq := graphql.RawParams{
				Query: `
mutation createMessage($input: NewMessage!) {
 createMessage(input: $input) {
   id
   message
   user {
	 id
	 name
   }
 }
}
		`,
				Variables: map[string]interface{}{
					"input": map[string]interface{}{
						"userID":  "AAA",
						"message": "new message 1",
					},
				},
			}
			req, _ := json.Marshal(graphReq)

			res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
			assert.NoError(t, err)
			// read body
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			assert.NoError(t, err)
			assert.NotEmpty(t, body)

			var gqlRes graphql.Response
			err = json.Unmarshal(body, &gqlRes)
			assert.NoError(t, err)
			assert.NotEmpty(t, gqlRes.Errors)
			assert.NotEmpty(t, gqlRes.Data)

			var data interface{}
			err = json.Unmarshal(gqlRes.Data, &data)
			assert.Nil(t, data)
		})

		t.Run("invalid params (lack param)", func(t *testing.T) {
			graphReq := graphql.RawParams{
				Query: `
mutation createMessage($input: NewMessage!) {
 createMessage(input: $input) {
   id
   message
   user {
	 id
	 name
   }
 }
}
		`,
				Variables: map[string]interface{}{
					"input": map[string]interface{}{
						"userID": 1,
					},
				},
			}
			req, _ := json.Marshal(graphReq)

			res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
			assert.NoError(t, err)
			// read body
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			assert.NoError(t, err)
			assert.NotEmpty(t, body)

			var gqlRes graphql.Response
			err = json.Unmarshal(body, &gqlRes)
			assert.NoError(t, err)
			assert.NotEmpty(t, gqlRes.Errors)
			assert.NotEmpty(t, gqlRes.Data)

			var data interface{}
			err = json.Unmarshal(gqlRes.Data, &data)
			assert.Nil(t, data)
		})

		t.Run("undefined schema", func(t *testing.T) {
			graphReq := graphql.RawParams{
				Query: `
mutation getWrong($input: NewMessage!) {
 wrong
}
`,
			}
			req, _ := json.Marshal(graphReq)

			res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
			assert.NoError(t, err)
			// read body
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			assert.NoError(t, err)
			assert.NotEmpty(t, body)

			var gqlRes graphql.Response
			err = json.Unmarshal(body, &gqlRes)
			assert.NoError(t, err)
			assert.NotNil(t, gqlRes.Errors)
			assert.NotEmpty(t, gqlRes.Errors)
			assert.NotNil(t, gqlRes.Data)

			var data interface{}
			err = json.Unmarshal(gqlRes.Data, &data)
			assert.NoError(t, err)
			assert.Nil(t, data)
		})
	})
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
