package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hilakatz/library/handlers"
)

func SetRoutes(incomingRoutes *gin.Engine, handler handlers.Handler) {
	books := incomingRoutes.Group("/books")
	{
		books.PUT("", handler.PutNewBook)
		books.POST("", handler.PostBookName)
		books.GET("", handler.GetBook)
		books.DELETE("", handler.DeleteBook)
	}

	search := incomingRoutes.Group("/search")
	{
		search.GET("", handler.SearchBooks)
	}

	store := incomingRoutes.Group("/store")
	{
		store.GET("", handler.GetInventory)
	}
}
