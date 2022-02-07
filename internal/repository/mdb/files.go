package mdb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"rendezvous/internal/model"
)

type FileRepository struct {
	c *mongo.Collection
}

const (
	fileCollection = "files"
)

func NewFilesRepository(db *mongo.Database) *UsersRepository {
	return &UsersRepository{c: db.Collection(fileCollection)}
}

func (f *FileRepository) CreateOrUpdateFiles(ctx context.Context, record *model.FileRecord) error {
	upsert := true
	if _, err := f.c.UpdateOne(ctx, bson.D{}, record, &options.UpdateOptions{Upsert: &upsert}); err != nil {
		return fmt.Errorf("failure to update meta files: %s", err)
	}
	return nil
}
