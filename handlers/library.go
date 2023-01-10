package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	consts "github.com/hilakatz/library/config"
	"github.com/hilakatz/library/models"
	"github.com/hilakatz/library/service"
	"net/http"
	"strings"
)

var validate = validator.New()

func PutNewBook(c *gin.Context) {

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

	result, insertErr := service.AddBook(*book.Title, *book.AuthorName, *book.Price, book.EbookAvailable, book.PublishDate)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Book item was not created."})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func PostBookName(c *gin.Context) {
	idString := c.Query(consts.ID)
	title := c.Query(consts.Title)

	result, updateErr := service.ChangeName(idString, title)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Book item was not updated."})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No matched book"})
		return
	}
	c.JSON(http.StatusCreated, idString)
}

func GetBook(c *gin.Context) {
	idString := c.Query(consts.ID)
	findErr, result := service.FindBook(idString)
	if findErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Book item was not found."})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func DeleteBook(c *gin.Context) {
	idString := c.Query(consts.ID)

	result, deleteErr := service.DeleteBook(idString)
	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Book item was not deleted."})
		return
	}
	c.JSON(http.StatusOK, result)
}

func SearchBooks(c *gin.Context) {

	title := c.Query(consts.Title)
	authorName := c.Query(consts.AuthorName)
	priceRange := c.Query(consts.PriceRange)
	priceRangeValues := strings.Split(priceRange, "-")

	books, findErr := service.FindBooksByParams(title, authorName, priceRangeValues)
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

func GetInventory(c *gin.Context) {

	authorResult, booksResult, err := service.RetrieveStore()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Distinct authors": len(authorResult), "Books": booksResult})
}
