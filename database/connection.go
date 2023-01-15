package database

import (
	"context"
	"log"

	errors "github.com/fiverr/go_errors"

	consts "github.com/hilakatz/library/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(consts.MongodbUrl))
	if err != nil {
		return nil, errors.Wrap(err, "failed to creat new client in Mongo")
	}

	ctx, cancel := context.WithTimeout(context.Background(), consts.MongoTimeout)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to connect to Mongo")
	}
	//ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, errors.Wrap(err, "failed to connect to Mongo")
	}

	log.Println("Connected to MongoDB")

	return client, nil
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database(consts.DbName).Collection(collectionName)

	return collection
}
