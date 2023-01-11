package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hilakatz/library/consts"
	"github.com/hilakatz/library/models"
	"net/http"
	"strings"
)

var validate = validator.New()

type Handler struct {
	library models.Library
}

func NewHandler(library models.Library) Handler {
	return Handler{library: library}
}

func (handler Handler) PutNewBook(c *gin.Context) {

	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := validate.Struct(book)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	insertErr := handler.library.AddBook(*book.Title, *book.AuthorName, *book.Price, book.EbookAvailable, book.PublishDate)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Book item was not created."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Book item was created."})

}

func (handler Handler) PostBookName(c *gin.Context) {
	idString := c.Query(consts.ID)
	title := c.Query(consts.Title)

	numUpdated, updateErr := handler.library.ChangeName(idString, title)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Book item was not updated."})
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
	findErr, result := handler.library.FindBook(idString)
	if findErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Book item was not found."})
		return
	}
	var data map[string]interface{}
	if err := json.Unmarshal(result, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, data)
}

func (handler Handler) DeleteBook(c *gin.Context) {
	idString := c.Query(consts.ID)

	deleteErr := handler.library.DeleteBook(idString)
	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Book item was not deleted."})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": findErr.Error()})
		return
	}

	if books == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no books where found"})
		return
	}

	c.JSON(http.StatusCreated, books)
}

func (handler Handler) GetInventory(c *gin.Context) {

	authorResult, booksResult, err := handler.library.RetrieveStore()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Distinct authors": len(authorResult), "Books": booksResult})
}
