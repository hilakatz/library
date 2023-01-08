package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hilakatz/library/consts"
	"github.com/hilakatz/library/controller"
	"github.com/hilakatz/library/routes"
)

func main() {

	port := consts.PORT
	router := gin.New()

	router.Use(controller.RequestLogger())

	routes.BookRoutes(router)
	routes.SearchRoutes(router)
	routes.StoreRoutes(router)

	router.Run("localhost:" + port)

}
