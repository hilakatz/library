package main

import (
	"fmt"

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
	mongo := service.NewMongoLibrary(setup.Client)
	handler := handlers.NewHandler(mongo)

	router := gin.New()
	router.Use(middlewares.RequestLogger, gin.Logger())
	routes.SetRoutes(router, handler)
	if err := router.Run(fmt.Sprintf("localhost:%s", config.PORT)); err != nil {
		panic(err)
	}
}
