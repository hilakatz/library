package service

import (
	"context"
	"log"
	"strconv"
	"time"

	consts "github.com/hilakatz/library/config"
	queryparams "github.com/hilakatz/library/consts"
	"github.com/hilakatz/library/database"
	"github.com/hilakatz/library/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoLibrary struct {
	bookCollection *mongo.Collection
}

func NewMongoLibrary(client *mongo.Client) *MongoLibrary {
	return &MongoLibrary{
		bookCollection: database.OpenCollection(client, consts.CollectionName),
	}
}

func (m MongoLibrary) AddBook(title, authorName string, price float64, ebookAvailable bool, publishDate time.Time) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	var book models.Book
	{
		book.ID = primitive.NewObjectID()
		book.Title = title
		book.AuthorName = authorName
		book.Price = price
		book.PublishDate = publishDate
		book.EbookAvailable = ebookAvailable
	}
	_, insertErr := m.bookCollection.InsertOne(ctx, book)
	return book.ID.Hex(), insertErr
}

func (m MongoLibrary) ChangeName(idString, title string) (int, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	opts := options.Update().SetUpsert(false)
	docID, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return 0, err
	}

	filter := bson.D{{queryparams.ID, docID}}
	update := bson.D{{"$set", bson.D{{queryparams.Title, title}}}}
	result, updateErr := m.bookCollection.UpdateOne(ctx, filter, update, opts)
	if result.MatchedCount == 0 {
		return 0, updateErr
	}

	return 1, updateErr
}

func (m MongoLibrary) FindBook(idString string) (models.Book, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	var book models.Book
	docID, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return models.Book{}, err
	}

	filter := bson.D{{queryparams.ID, docID}}
	findErr := m.bookCollection.FindOne(ctx, filter).Decode(&book)
	if findErr != nil {
		return models.Book{}, findErr
	}

	return book, nil
}

func (m MongoLibrary) DeleteBook(idString string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	docID, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return err
	}

	filter := bson.D{{queryparams.ID, docID}}
	_, deleteErr := m.bookCollection.DeleteOne(ctx, filter)
	return deleteErr
}

func (m MongoLibrary) FindBooksByParams(title, authorName string, priceRangeValues []string) ([]models.Book, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	filter := bson.M{}
	if title != "" {
		filter[queryparams.Title] = title
	}
	if authorName != "" {
		filter[queryparams.AuthorName] = authorName
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
		filter[queryparams.Price] = bson.M{"$gte": priceMin, "$lte": priceMax}
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

func (m MongoLibrary) RetrieveStore() (int, int, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	filter := bson.D{{}}
	authorResult, err := m.bookCollection.Distinct(ctx, queryparams.AuthorName, filter)
	if err != nil {
		return 0, 0, err
	}

	booksResult, err := m.bookCollection.EstimatedDocumentCount(ctx)
	if err != nil {
		return 0, 0, err
	}

	return len(authorResult), int(booksResult), nil
}
