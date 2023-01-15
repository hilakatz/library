package setup

import (
	"github.com/hilakatz/library/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client

func Setup() error {
	var err error
	MongoClient, err = database.ConnectDB()

	return err
}
