package persistence

import (
	"context"
	"database/sql"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/domain/entity"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/graph/model"
	"github.com/google/go-cmp/cmp"
	"github.com/graph-gophers/dataloader"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepo_ListUsers(t *testing.T) {
	t.Run("get all users", func(t *testing.T) {
		users, err := userRepo.ListUsers()

		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.NotEmpty(t, users)
		assert.Len(t, users, 2)

		want := []entity.User{
			{Id: 2, Name: "Fuga"},
			{Id: 1, Name: "Hoge"},
		}
		if diff := cmp.Diff(want, users); len(diff) != 0 {
			t.Errorf("Users mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestUserRepo_User(t *testing.T) {
	t.Run("get user=1", func(t *testing.T) {
		user, err := userRepo.User(1)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotEmpty(t, user)
		assert.Equal(t, 1, user.Id)
		assert.Equal(t, "Hoge", user.Name)
	})
	t.Run("get user=2", func(t *testing.T) {
		user, err := userRepo.User(2)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotEmpty(t, user)
		assert.Equal(t, 2, user.Id)
		assert.Equal(t, "Fuga", user.Name)
	})
	t.Run("get not exist user", func(t *testing.T) {
		user, err := userRepo.User(9999)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
		assert.NotNil(t, user)
	})
}

func TestUserRepo_GetUsers(t *testing.T) {
	t.Run("get user=1", func(t *testing.T) {
		ctx := context.Background()
		keys := dataloader.NewKeysFromStrings([]string{"1"})
		result := userRepo.GetUsers(ctx, keys)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result)
		assert.Len(t, result, 1)

		assert.Nil(t, result[0].Error)
		user := result[0].Data.(*model.User)
		assert.NotNil(t, user)
		assert.Equal(t, "1", user.ID)
		assert.Equal(t, "Hoge", user.Name)
	})
	t.Run("get user=2", func(t *testing.T) {
		ctx := context.Background()
		keys := dataloader.NewKeysFromStrings([]string{"2"})
		result := userRepo.GetUsers(ctx, keys)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result)
		assert.Len(t, result, 1)

		assert.Nil(t, result[0].Error)
		user := result[0].Data.(*model.User)
		assert.NotNil(t, user)
		assert.Equal(t, "2", user.ID)
		assert.Equal(t, "Fuga", user.Name)
	})

	t.Run("get user=1, 2", func(t *testing.T) {
		ctx := context.Background()
		keys := dataloader.NewKeysFromStrings([]string{"1", "2"})
		result := userRepo.GetUsers(ctx, keys)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result)
		assert.Len(t, result, 2)

		assert.Nil(t, result[0].Error)
		user1 := result[0].Data.(*model.User)
		assert.NotNil(t, user1)
		assert.Equal(t, "1", user1.ID)
		assert.Equal(t, "Hoge", user1.Name)

		assert.Nil(t, result[1].Error)
		user2 := result[1].Data.(*model.User)
		assert.NotNil(t, user2)
		assert.Equal(t, "2", user2.ID)
		assert.Equal(t, "Fuga", user2.Name)
	})

	t.Run("get not exist user", func(t *testing.T) {
		ctx := context.Background()
		keys := dataloader.NewKeysFromStrings([]string{"9999"})
		result := userRepo.GetUsers(ctx, keys)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result)
		assert.Len(t, result, 1)
		assert.NotNil(t, result[0].Error)
	})
	t.Run("get user=1, not exist user", func(t *testing.T) {
		ctx := context.Background()
		keys := dataloader.NewKeysFromStrings([]string{"1", "9999"})
		result := userRepo.GetUsers(ctx, keys)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result)
		assert.Len(t, result, 2)

		assert.Nil(t, result[0].Error)
		user := result[0].Data.(*model.User)
		assert.NotNil(t, user)
		assert.Equal(t, "1", user.ID)
		assert.Equal(t, "Hoge", user.Name)

		assert.NotNil(t, result[1].Error)
	})

}
