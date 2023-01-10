package config

import "time"

// MONGO
const (
	MongodbUrl     = "mongodb://localhost:27017"
	MongoTimeout   = 10 * time.Second
	HandlerTimeout = 100 * time.Second
	DbName         = "library"
	CollectionName = "books"
)

// PORT SERVER
const (
	PORT = "8080"
)
