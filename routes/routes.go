package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hilakatz/library/handlers"
)

func SetRoutes(incomingRoutes *gin.Engine, handler handlers.Handler) {
	api := incomingRoutes.Group("/api")
	{
		books := api.Group("/books")
		{
			books.PUT("/", handler.PutNewBook)
			books.POST("/:id", handler.PostBookName)
			books.GET("/:id", handler.GetBook)
			books.DELETE("/:id", handler.DeleteBook)
		}

		search := api.Group("/search")
		{
			search.GET("/", handler.SearchBooks)
		}

		store := api.Group("/store")
		{
			store.GET("/", handler.GetInventory)
		}

	}
}
