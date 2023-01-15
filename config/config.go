package config

import "time"

const (
	MongodbUrl     = "mongodb://localhost:27017"
	MongoTimeout   = 10 * time.Second
	HandlerTimeout = 100 * time.Second
	DbName         = "library"
	CollectionName = "books"
	PORT           = "8080"
)
