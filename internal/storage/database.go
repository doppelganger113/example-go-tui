package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewDbConnection(ctx context.Context, uri string) (*mongo.Client, error) {
	var updatedUri string
	if uri == "" {
		updatedUri = "mongodb://localhost:27017"
	}

	return mongo.Connect(ctx, options.Client().
		ApplyURI(updatedUri).
		SetTimeout(5*time.Second),
	)
}
