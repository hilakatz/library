package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hilakatz/library/handlers"
)

func GetRoutes(incomingRoutes *gin.Engine) {
	books := incomingRoutes.Group("/books")
	{
		books.PUT("", handlers.PutNewBook)
		books.POST("", handlers.PostBookName)
		books.GET("", handlers.GetBook)
		books.DELETE("", handlers.DeleteBook)
	}
	search := incomingRoutes.Group("/search")
	{
		search.GET("", handlers.SearchBooks)
	}
	store := incomingRoutes.Group("/store")
	{
		store.GET("", handlers.GetInventory)
	}
}
