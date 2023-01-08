package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hilakatz/library/database"
	"github.com/hilakatz/library/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var bookCollection *mongo.Collection = database.OpenCollection(database.Client, "books")

var validate = validator.New()

func PutNewBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var book model.Book

		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(book)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		book.ID = primitive.NewObjectID()

		result, insertErr := bookCollection.InsertOne(ctx, book)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Book item was not created."})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func PostBookName() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		opts := options.Update().SetUpsert(false)
		idString := c.Query("_id")
		docID, err := primitive.ObjectIDFromHex(idString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		filter := bson.D{{"_id", docID}}
		title := c.Query("title")
		update := bson.D{{"$set", bson.D{{"title", title}}}}

		result, updateErr := bookCollection.UpdateOne(ctx, filter, update, opts)
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Book item was not updated."})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No matched book"})
			return
		}
		c.JSON(http.StatusCreated, docID)
	}
}

func GetBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var result bson.M
		idString := c.Query("_id")
		docID, err := primitive.ObjectIDFromHex(idString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		filter := bson.D{{"_id", docID}}

		findErr := bookCollection.FindOne(ctx, filter).Decode(&result)
		if findErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Book item was not found."})
			return
		}
		c.JSON(http.StatusCreated, result)
	}
}

func DeleteBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		idString := c.Query("_id")
		docID, err := primitive.ObjectIDFromHex(idString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		filter := bson.D{{"_id", docID}}

		result, deleteErr := bookCollection.DeleteOne(ctx, filter)
		if deleteErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Book item was not deleted."})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func SearchBooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		title := c.Query("title")
		authorName := c.Query("author_name")
		priceRange := c.Query("price_range")
		priceRangeValues := strings.Split(priceRange, "-")

		filter := bson.M{}
		if title != "" {
			filter["title"] = title
		}
		if authorName != "" {
			filter["author_name"] = authorName
		}

		if len(priceRangeValues) == 2 {
			priceMin, err := strconv.Atoi(priceRangeValues[0])
			if err != nil {
				priceMin = 0
			}

			priceMax, err := strconv.Atoi(priceRangeValues[1])
			if err != nil {
				priceMax = 1000000
			}

			filter["price"] = bson.M{"$gte": priceMin, "$lte": priceMax}
		}
		cur, findErr := bookCollection.Find(ctx, filter)
		if findErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": findErr.Error()})
			return
		}
		var books []model.Book
		if findErr = cur.All(ctx, &books); findErr != nil {
			log.Fatal(findErr)
		}

		if books == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "no books where found"})
			return
		}

		c.JSON(http.StatusCreated, books)
	}
}

func GetInventory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{{}}
		authorResult, err := bookCollection.Distinct(ctx, "author_name", filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		booksResult, err := bookCollection.EstimatedDocumentCount(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"Distinct authors": len(authorResult), "Books": booksResult})
	}
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params []string
		for _, val := range c.Request.URL.Query() {
			for _, v := range val {
				params = append(params, v)
			}
		}
		// Print the request details to the log
		fmt.Printf("%s - %s::%s::%v\n",
			time.Now().Format(time.RFC3339),
			c.Request.Method,
			c.Request.URL.Path,
			strings.Join(params, `,`),
		)
		// Continue with the request
		c.Next()
	}
}
