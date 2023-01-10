package service

import (
	"context"
	consts "github.com/hilakatz/library/config"
	"github.com/hilakatz/library/database"
	"github.com/hilakatz/library/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"time"
)

var bookCollection *mongo.Collection = database.OpenCollection(database.Client, consts.CollectionName)

func AddBook(title, authorName string, price float64, ebookAvailable bool, publishDate time.Time) (*mongo.InsertOneResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	var book models.Book
	{
		book.ID = primitive.NewObjectID()
		book.Title = &title
		book.AuthorName = &authorName
		book.Price = &price
		book.PublishDate = publishDate
		book.EbookAvailable = ebookAvailable
	}
	result, insertErr := bookCollection.InsertOne(ctx, book)
	return result, insertErr
}

func ChangeName(idString, title string) (*mongo.UpdateResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()
	opts := options.Update().SetUpsert(false)
	docID, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", docID}}
	update := bson.D{{"$set", bson.D{{"title", title}}}}
	result, updateErr := bookCollection.UpdateOne(ctx, filter, update, opts)
	return result, updateErr
}

func FindBook(idString string) (error, bson.M) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	var result bson.M
	docID, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return err, nil
	}
	filter := bson.D{{"_id", docID}}

	findErr := bookCollection.FindOne(ctx, filter).Decode(&result)
	return findErr, result
}

func DeleteBook(idString string) (*mongo.DeleteResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	docID, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", docID}}
	result, deleteErr := bookCollection.DeleteOne(ctx, filter)
	return result, deleteErr
}

func FindBooksByParams(title, authorName string, priceRangeValues []string) ([]models.Book, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

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
		return nil, findErr
	}
	var books []models.Book
	if findErr = cur.All(ctx, &books); findErr != nil {
		log.Fatal(findErr)
	}

	return books, nil

}

func RetrieveStore() ([]interface{}, int64, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	filter := bson.D{{}}
	authorResult, err := bookCollection.Distinct(ctx, "author_name", filter)
	if err != nil {
		return nil, 0, err
	}
	booksResult, err := bookCollection.EstimatedDocumentCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return authorResult, booksResult, nil
}
