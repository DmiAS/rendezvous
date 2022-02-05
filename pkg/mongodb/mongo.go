package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const timeout = 1 * time.Second

func NewConnection(dsn, database string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, fmt.Errorf("failure to connect to database: %s", err)
	}

	// check the connection
	if err := conn.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failture to check connection: %s", err)
	}
	return conn.Database(database), nil
}
