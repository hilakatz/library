package service

import (
	"context"
	"fmt"
	"log"
	"time"

	errors "github.com/fiverr/go_errors"

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
		book.Title = title
		book.AuthorName = authorName
		book.Price = price
		book.PublishDate = publishDate
		book.EbookAvailable = ebookAvailable
	}

	InsertOneResult, insertErr := m.bookCollection.InsertOne(ctx, book)
	if insertErr != nil {
		return "", errors.Wrap(insertErr, "failed to insert to Mongo")
	}

	return fmt.Sprintf("%v", InsertOneResult.InsertedID), nil
}

func (m MongoLibrary) ChangeName(idString, title string) (int, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	opts := options.Update().SetUpsert(false)
	docID, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return 0, errors.Wrap(err, "failed to convert - not valid ObjectID")
	}

	filter := bson.D{{queryparams.ID, docID}}
	update := bson.D{{"$set", bson.D{{queryparams.Title, title}}}}
	result, updateErr := m.bookCollection.UpdateOne(ctx, filter, update, opts)
	if result.MatchedCount == 0 {
		return 0, errors.Wrap(updateErr, "failed to update Mongo")
	}

	return 1, nil
}

func (m MongoLibrary) FindBook(idString string) (models.Book, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	var book models.Book
	docID, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return models.Book{}, errors.Wrap(err, "failed to convert - not valid ObjectID")
	}

	filter := bson.D{{queryparams.ID, docID}}
	findErr := m.bookCollection.FindOne(ctx, filter).Decode(&book)
	if findErr != nil {
		return models.Book{}, errors.Wrap(findErr, "failed to find Document in Mongo")
	}

	return book, nil
}

func (m MongoLibrary) DeleteBook(idString string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	docID, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return errors.Wrap(err, "failed to convert - not valid ObjectID")
	}

	filter := bson.D{{queryparams.ID, docID}}
	result, deleteErr := m.bookCollection.DeleteOne(ctx, filter)
	if result.DeletedCount == 0 {
		return errors.Wrap(deleteErr, "failed to delete document from Mongo")
	}
	return nil
}

func (m MongoLibrary) FindBooksByParams(title, authorName string, priceMin, priceMax float64) ([]models.Book, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	filter := bson.M{}
	if title != "" {
		filter[queryparams.Title] = title
	}
	if authorName != "" {
		filter[queryparams.AuthorName] = authorName
	}
	filter[queryparams.Price] = bson.M{"$gte": priceMin, "$lte": priceMax}

	cur, findErr := m.bookCollection.Find(ctx, filter)
	if findErr != nil {
		return nil, errors.Wrap(findErr, "failed to find Documents in Mongo")
	}

	var books []models.Book
	if findErr = cur.All(ctx, &books); findErr != nil {
		log.Fatal(errors.Wrap(findErr, "failed to find Documents in Mongo"))
	}

	return books, nil
}

func (m MongoLibrary) RetrieveStore() (int, int, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), consts.HandlerTimeout)
	defer cancel()

	filter := bson.D{{}}
	authorResult, err := m.bookCollection.Distinct(ctx, queryparams.AuthorName, filter)
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to find Documents in Mongo")
	}

	booksResult, err := m.bookCollection.EstimatedDocumentCount(ctx)
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to find Documents count in Mongo")
	}

	return len(authorResult), int(booksResult), nil
}
