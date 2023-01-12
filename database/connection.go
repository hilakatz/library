package database

import (
	"context"
	"log"

	consts "github.com/hilakatz/library/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(consts.MongodbUrl))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), consts.MongoTimeout)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}
	//ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB")

	return client, nil
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database(consts.DbName).Collection(collectionName)

	return collection
}
