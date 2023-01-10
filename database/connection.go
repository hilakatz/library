package database

import (
	"context"
	consts "github.com/hilakatz/library/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(consts.MongodbUrl))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), consts.MongoTimeout)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		panic(err)
	}

	//ping the database
	if err := client.Ping(ctx, nil); err != nil {
		panic(err)
	}
	log.Println("Connected to MongoDB")
	return client
}

var Client *mongo.Client = ConnectDB()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database(consts.DbName).Collection(collectionName)
	return collection
}
