package main

import (
	"fmt"

	errors "github.com/fiverr/go_errors"

	"github.com/hilakatz/library/setup"

	"github.com/gin-gonic/gin"
	"github.com/hilakatz/library/config"
	"github.com/hilakatz/library/handlers"
	"github.com/hilakatz/library/middlewares"
	"github.com/hilakatz/library/routes"
	"github.com/hilakatz/library/service"
)

func main() {
	if err := setup.Setup(); err != nil {
		panic(err)
	}

	mongo := service.NewMongoLibrary(setup.MongoClient)
	handler := handlers.NewHandler(mongo)

	router := gin.New()
	router.Use(middlewares.RequestLogger, gin.Logger())
	routes.SetRoutes(router, handler)
	if err := router.Run(fmt.Sprintf("localhost:%s", config.PORT)); err != nil {
		panic(errors.Wrap(err, "failed to attache the router to server"))
	}
}
