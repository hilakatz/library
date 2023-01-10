package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hilakatz/library/config"
	"github.com/hilakatz/library/handlers"
	"github.com/hilakatz/library/routes"
)

func main() {
	router := gin.New()

	router.Use(handlers.RequestLogger)

	routes.GetRoutes(router)

	if err := router.Run(fmt.Sprintf("localhost:%s", config.PORT)); err != nil {
		panic(err)
	}
}
