package persistence

import (
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/domain/entity"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessageRepo_ListMessages(t *testing.T) {

	t.Run("get messages user=1", func(t *testing.T) {
		messages, err := messageRepo.ListMessages(1)
		assert.NoError(t, err)
		assert.NotNil(t, messages)
		assert.NotEmpty(t, messages)
		assert.Len(t, messages, 2)

		want := []entity.Message{
			{Id: 2, UserId: 1, Message: "test message id 2"},
			{Id: 1, UserId: 1, Message: "test message id 1"},
		}
		if diff := cmp.Diff(want, messages); len(diff) != 0 {
			t.Errorf("Messages mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("get messages user=2", func(t *testing.T) {
		messages, err := messageRepo.ListMessages(2)
		assert.NoError(t, err)
		assert.NotNil(t, messages)
		assert.NotEmpty(t, messages)
		assert.Len(t, messages, 2)

		want := []entity.Message{
			{Id: 4, UserId: 2, Message: "test message id 4"},
			{Id: 3, UserId: 2, Message: "test message id 3"},
		}
		if diff := cmp.Diff(want, messages); len(diff) != 0 {
			t.Errorf("Messages mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("get messages not exist user", func(t *testing.T) {
		messages, err := messageRepo.ListMessages(99999)
		assert.NoError(t, err)
		assert.NotNil(t, messages)
		assert.Empty(t, messages)
	})

}

func TestMessageRepo_Messages(t *testing.T) {
	t.Run("get messages", func(t *testing.T) {
		messages, err := messageRepo.Messages()
		assert.NoError(t, err)
		assert.NotNil(t, messages)
		assert.NotEmpty(t, messages)
		assert.Len(t, messages, 4)

		want := []entity.Message{
			{Id: 4, UserId: 2, Message: "test message id 4"},
			{Id: 3, UserId: 2, Message: "test message id 3"},
			{Id: 2, UserId: 1, Message: "test message id 2"},
			{Id: 1, UserId: 1, Message: "test message id 1"},
		}
		if diff := cmp.Diff(want, messages); len(diff) != 0 {
			t.Errorf("Messages mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("get blank ", func(t *testing.T) {
		org, err := messageRepo.Messages()
		db := getDatabase()

		stmt, _ := db.Prepare("TRUNCATE TABLE message;")
		stmt.Exec()
		stmt.Close()

		messages, err := messageRepo.Messages()
		assert.NoError(t, err)
		assert.NotNil(t, messages)
		assert.Empty(t, messages)
		assert.Len(t, messages, 0)

		t.Cleanup(func() {
			stmt, _ = db.Prepare("INSERT INTO message(id, user_id, message) VALUES(?, ?, ?);")
			for _, msg := range org {
				stmt.Exec(msg.Id, msg.UserId, msg.Message)
			}
			stmt.Close()
		})
	})
}

func TestMessageRepo_CreateMessage(t *testing.T) {
	t.Run("create messages user=1", func(t *testing.T) {
		db := getDatabase()

		message := &entity.Message{
			UserId:  1,
			Message: "new message 1",
		}
		err := messageRepo.CreateMessage(message)
		assert.NoError(t, err)
		assert.NotNil(t, message)
		assert.NotEmpty(t, message)

		assert.Greater(t, message.Id, 1)

		stmt, err := db.Prepare("SELECT id, user_id, message FROM message WHERE id = ?")
		defer stmt.Close()
		if err != nil {
			panic(err.Error())
		}
		newMessage := entity.Message{}
		err = stmt.QueryRow(message.Id).Scan(&newMessage.Id, &newMessage.UserId, &newMessage.Message)
		if err != nil {
			panic(err.Error())
		}
		assert.NotEmpty(t, newMessage)
		assert.Equal(t, message.UserId, newMessage.UserId)
		assert.Equal(t, message.Message, newMessage.Message)

		t.Cleanup(func() {
			if message != nil && message.Id > 0 {
				st, er := db.Prepare("DELETE FROM message WHERE id = ?")
				if er == nil {
					_, er = st.Exec(message.Id)
					if er != nil {
						panic(er.Error())
					}
				}
				defer st.Close()
			}
		})
	})

	t.Run("create messages user=2", func(t *testing.T) {
		db := getDatabase()

		message := &entity.Message{
			UserId:  2,
			Message: "new message 2",
		}
		err := messageRepo.CreateMessage(message)
		assert.NoError(t, err)
		assert.NotNil(t, message)
		assert.NotEmpty(t, message)

		assert.Greater(t, message.Id, 1)

		stmt, err := db.Prepare("SELECT id, user_id, message FROM message WHERE id = ?")
		defer stmt.Close()
		if err != nil {
			panic(err.Error())
		}
		newMessage := entity.Message{}
		err = stmt.QueryRow(message.Id).Scan(&newMessage.Id, &newMessage.UserId, &newMessage.Message)
		if err != nil {
			panic(err.Error())
		}
		assert.NotEmpty(t, newMessage)
		assert.Equal(t, message.UserId, newMessage.UserId)
		assert.Equal(t, message.Message, newMessage.Message)

		t.Cleanup(func() {
			if message != nil && message.Id > 0 {
				st, er := db.Prepare("DELETE FROM message WHERE id = ?")
				if er == nil {
					_, er = st.Exec(message.Id)
					if er != nil {
						panic(er.Error())
					}
				}
				defer st.Close()
			}
		})
	})

	t.Run("create not exist user_id", func(t *testing.T) {
		message := &entity.Message{
			UserId:  9999,
			Message: "new message 3",
		}
		err := messageRepo.CreateMessage(message)
		assert.Error(t, err)
		assert.Equal(t, message.Id, 0)
	})
}
