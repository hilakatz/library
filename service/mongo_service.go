package service

import (
	"context"
	"encoding/json"
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

type MongoLibrary struct {
	bookCollection *mongo.Collection
}

func NewMongoLibrary() *MongoLibrary {
	return &MongoLibrary{
		bookCollection: database.OpenCollection(database.Client, consts.CollectionName),
	}
}

func (m MongoLibrary) AddBook(title, authorName string, price float64, ebookAvailable bool, publishDate time.Time) error {
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
	_, insertErr := m.bookCollection.InsertOne(ctx, book)
	return insertErr
}

func (m MongoLibrary) ChangeName(idString, title string) (int, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()
	opts := options.Update().SetUpsert(false)
	docID, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return 0, err
	}
	filter := bson.D{{"_id", docID}}
	update := bson.D{{"$set", bson.D{{"title", title}}}}
	result, updateErr := m.bookCollection.UpdateOne(ctx, filter, update, opts)
	if result.MatchedCount == 0 {
		return 0, updateErr
	}
	return 1, updateErr
}

func (m MongoLibrary) FindBook(idString string) (error, []byte) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	var result bson.M
	docID, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return err, nil
	}
	filter := bson.D{{"_id", docID}}

	findErr := m.bookCollection.FindOne(ctx, filter).Decode(&result)
	if findErr != nil {
		return findErr, nil
	}
	jsonResult, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		return jsonErr, nil
	}
	return nil, jsonResult
}

func (m MongoLibrary) DeleteBook(idString string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	docID, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", docID}}
	_, deleteErr := m.bookCollection.DeleteOne(ctx, filter)
	return deleteErr
}

func (m MongoLibrary) FindBooksByParams(title, authorName string, priceRangeValues []string) ([]models.Book, error) {
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
	cur, findErr := m.bookCollection.Find(ctx, filter)
	if findErr != nil {
		return nil, findErr
	}
	var books []models.Book
	if findErr = cur.All(ctx, &books); findErr != nil {
		log.Fatal(findErr)
	}

	return books, nil

}

func (m MongoLibrary) RetrieveStore() ([]interface{}, int64, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	filter := bson.D{{}}
	authorResult, err := m.bookCollection.Distinct(ctx, "author_name", filter)
	if err != nil {
		return nil, 0, err
	}
	booksResult, err := m.bookCollection.EstimatedDocumentCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return authorResult, booksResult, nil
}
