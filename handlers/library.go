package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hilakatz/library/consts"
	queryparams "github.com/hilakatz/library/consts"
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

	ID, err := handler.library.AddBook(book.Title, book.AuthorName, book.Price, book.EbookAvailable, book.PublishDate)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": ID})
}

func (handler Handler) PostBookName(c *gin.Context) {
	idString := c.Param("id")
	title := c.Query(consts.Title)

	_, err := handler.library.ChangeName(idString, title)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Book item was updated"})
}

func (handler Handler) GetBook(c *gin.Context) {
	idString := c.Param("id")
	result, err := handler.library.FindBook(idString)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (handler Handler) DeleteBook(c *gin.Context) {
	idString := c.Param("id")

	err := handler.library.DeleteBook(idString)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book item was deleted."})
}

func (handler Handler) SearchBooks(c *gin.Context) {
	var priceMin, priceMax float64
	title := c.Query(consts.Title)
	authorName := c.Query(consts.AuthorName)
	priceRange := c.Query(consts.PriceRange)
	priceRangeValues := strings.Split(priceRange, "-")
	if len(priceRangeValues) == 2 {
		var err error
		priceMin, err = strconv.ParseFloat(priceRangeValues[0], 64)
		if err != nil {
			priceMin = queryparams.PriceMin
		}
		priceMax, err = strconv.ParseFloat(priceRangeValues[1], 64)
		if err != nil {
			priceMax = queryparams.PriceMax
		}
	}

	books, err := handler.library.FindBooksByParams(title, authorName, priceMin, priceMax)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if books == nil {
		c.AbortWithError(http.StatusNotFound, err)
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
