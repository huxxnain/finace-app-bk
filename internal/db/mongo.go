package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Connect(mongoDBURI string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoDBURI))
	if err != nil {
		return err
	}

	// Verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func GetClient() *mongo.Client {
	return client
}

func GetDatabase(dbName string) *mongo.Database {
	return client.Database(dbName)
}

func Close() error {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return client.Disconnect(ctx)
	}
	return nil
}
