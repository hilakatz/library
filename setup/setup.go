package setup

import (
	"github.com/hilakatz/library/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var Client *mongo.Client

func Setup() error {
	var err error
	Client, err = database.ConnectDB()

	return err
}
