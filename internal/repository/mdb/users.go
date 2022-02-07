package mdb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"rendezvous/internal/model"
)

type UsersRepository struct {
	c *mongo.Collection
}

const (
	userCollection = "Users"
)

func NewUsersRepository(db *mongo.Database) *UsersRepository {
	return &UsersRepository{c: db.Collection(userCollection)}
}

func (u UsersRepository) GetUsers(ctx context.Context) (model.Users, error) {
	cursor, err := u.c.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failure to extract cursor: %s", err)
	}
	var users model.Users
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failure to get users: %s", err)
	}
	return users, nil
}

func (u UsersRepository) GetUser(ctx context.Context, login string) (*model.User, error) {
	cursor := u.c.FindOne(ctx, bson.D{{"name", login}})
	user := &model.User{}
	if err := cursor.Decode(user); err != nil {
		return nil, fmt.Errorf("failure to decode user %s: %s", login, err)
	}
	return user, nil
}

func (u UsersRepository) UpdateUser(ctx context.Context, user *model.User) error {
	updateData := map[string]interface{}{
		"chatting": user.Chatting,
	}
	if user.LocalAddress != "" {
		updateData["local_address"] = user.LocalAddress
	}
	if user.GlobalAddress != "" {
		updateData["global_address"] = user.GlobalAddress
	}
	_, err := u.c.UpdateOne(ctx, bson.D{{"name", user.Name}}, updateData)
	if err != nil {
		return fmt.Errorf("failure to update user %+v: %s", user, err)
	}
	return nil
}

func (u UsersRepository) CreateUser(ctx context.Context, user *model.User) error {
	if _, err := u.c.InsertOne(ctx, user); err != nil {
		return fmt.Errorf("failure to insert user: %s", err)
	}
	return nil
}

func (u UsersRepository) DeleteUser(ctx context.Context, login string) error {
	if _, err := u.c.DeleteOne(ctx, bson.D{{"login", login}}); err != nil {
		return fmt.Errorf("failure to delete user: %s", err)
	}
	return nil
}
