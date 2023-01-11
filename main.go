package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hilakatz/library/config"
	"github.com/hilakatz/library/handlers"
	"github.com/hilakatz/library/middlewares"
	"github.com/hilakatz/library/routes"
	"github.com/hilakatz/library/service"
)

func main() {
	router := gin.New()
	mongo := service.NewMongoLibrary()
	handler := handlers.NewHandler(mongo)
	router.Use(middlewares.RequestLogger)
	routes.GetRoutes(router, handler)
	if err := router.Run(fmt.Sprintf("localhost:%s", config.PORT)); err != nil {
		panic(err)
	}
}
