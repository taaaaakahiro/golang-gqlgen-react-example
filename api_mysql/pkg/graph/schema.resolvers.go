package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/domain/entity"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/graph/generated"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/graph/model"
)

func (r *messageResolver) User(ctx context.Context, obj *model.Message) (*model.User, error) {
	//// MEMO: ユーザを要求があったときだけ取得できる
	// ref https://gqlgen.com/reference/dataloaders/
	return r.Loaders.GetUser(ctx, obj.UserID)
}

func (r *mutationResolver) CreateMessage(ctx context.Context, input model.NewMessage) (*model.Message, error) {
	userId, err := strconv.Atoi(input.UserID)
	if err != nil {
		return nil, err
	}
	_, err = r.Repo.User.User(userId)
	if err != nil {
		// not exist etc...
		return nil, errors.New("user error. " + err.Error())
	}

	entityMessage := &entity.Message{
		UserId:  userId,
		Message: input.Message,
	}
	err = r.Repo.Message.CreateMessage(entityMessage)
	if err != nil {
		return nil, err
	}
	result := &model.Message{
		Message: input.Message,
		ID:      strconv.Itoa(entityMessage.Id),
		UserID:  input.UserID,
	}
	return result, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	entities, err := r.Repo.User.ListUsers()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	users := make([]*model.User, 0)
	for _, u := range entities {
		users = append(users, &model.User{
			ID:   strconv.Itoa(u.Id),
			Name: u.Name,
		})
	}
	return users, nil
}

func (r *queryResolver) Messages(ctx context.Context, userID string) ([]*model.Message, error) {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}
	entities, err := r.Repo.Message.ListMessages(id)
	if err != nil {
		return nil, err
	}
	messages := make([]*model.Message, 0)
	for _, ent := range entities {
		messages = append(messages, &model.Message{
			ID:      strconv.Itoa(ent.Id),
			Message: ent.Message,
			UserID:  strconv.Itoa(ent.UserId),
		})
	}
	return messages, nil
}

func (r *queryResolver) AllMessages(ctx context.Context) ([]*model.Message, error) {
	entities, err := r.Repo.Message.Messages()
	if err != nil {
		return nil, err
	}
	messages := make([]*model.Message, 0)
	for _, ent := range entities {
		messages = append(messages, &model.Message{
			ID:      strconv.Itoa(ent.Id),
			Message: ent.Message,
			UserID:  strconv.Itoa(ent.UserId),
		})
	}
	return messages, nil
}

// Message returns generated.MessageResolver implementation.
func (r *Resolver) Message() generated.MessageResolver { return &messageResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type messageResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
