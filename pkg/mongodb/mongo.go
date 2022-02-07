package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const timeout = 1 * time.Second

type Connection struct {
	Conn *mongo.Client
	Db   *mongo.Database
}

func NewConnection(dsn, database string) (*Connection, error) {
	conn, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, fmt.Errorf("failure to connect to database: %s", err)
	}

	// check the connection
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := conn.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failture to check connection: %s", err)
	}
	return &Connection{Conn: conn, Db: conn.Database(database)}, nil
}
