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
	collection = "users"
)

func NewUsersRepository(db *mongo.Database) *UsersRepository {
	return &UsersRepository{c: db.Collection(collection)}
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
