package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/hilakatz/library/controller"
)

func BookRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.PUT("/books", controllers.PutNewBook())
	incomingRoutes.POST("/books", controllers.PostBookName())
	incomingRoutes.GET("/books", controllers.GetBook())
	incomingRoutes.DELETE("/books", controllers.DeleteBook())
}

func SearchRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/search", controllers.SearchBooks())
}

func StoreRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/store", controllers.GetInventory())
}
