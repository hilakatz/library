package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hilakatz/library/consts"
	"github.com/hilakatz/library/models"
)

type Handler struct {
	library models.Library
}

func NewHandler(library models.Library) Handler {
	return Handler{library: library}
}

func (handler Handler) PutNewBook(c *gin.Context) {
	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ID, insertErr := handler.library.AddBook(book.Title, book.AuthorName, book.Price, book.EbookAvailable, book.PublishDate)
	if insertErr != nil {
		c.AbortWithError(http.StatusInternalServerError, insertErr)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"_id": ID})
}

func (handler Handler) PostBookName(c *gin.Context) {
	idString := c.Query(consts.ID)
	title := c.Query(consts.Title)

	numUpdated, updateErr := handler.library.ChangeName(idString, title)
	if updateErr != nil {
		c.AbortWithError(http.StatusInternalServerError, updateErr)
		return
	}
	if numUpdated == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No matched book"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Book item was updated"})
}

func (handler Handler) GetBook(c *gin.Context) {
	idString := c.Query(consts.ID)

	result, findErr := handler.library.FindBook(idString)
	if findErr != nil {
		c.AbortWithError(http.StatusInternalServerError, findErr)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (handler Handler) DeleteBook(c *gin.Context) {
	idString := c.Query(consts.ID)

	deleteErr := handler.library.DeleteBook(idString)
	if deleteErr != nil {
		c.AbortWithError(http.StatusInternalServerError, deleteErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book item was deleted."})
}

func (handler Handler) SearchBooks(c *gin.Context) {
	title := c.Query(consts.Title)
	authorName := c.Query(consts.AuthorName)
	priceRange := c.Query(consts.PriceRange)
	priceRangeValues := strings.Split(priceRange, "-")

	books, findErr := handler.library.FindBooksByParams(title, authorName, priceRangeValues)
	if findErr != nil {
		c.AbortWithError(http.StatusInternalServerError, findErr)
		return
	}

	if books == nil {
		c.AbortWithError(http.StatusNotFound, findErr)
		return
	}

	c.JSON(http.StatusCreated, books)
}

func (handler Handler) GetInventory(c *gin.Context) {
	authorResult, booksResult, err := handler.library.RetrieveStore()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Distinct authors": authorResult, "Books": booksResult})
}
