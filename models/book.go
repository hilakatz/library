package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Book struct {
	ID             primitive.ObjectID `bson:"_id"`
	Title          *string            `json:"title" bson:"title" validate:"required"`
	AuthorName     *string            `json:"author_name" bson:"author_name" validate:"required"`
	Price          *float64           `json:"price" bson:"price" validate:"required"`
	EbookAvailable bool               `json:"ebook_available" bson:"ebook_available"`
	PublishDate    time.Time          `json:"publish_date" bson:"publish_date" validate:"required"`
}

type Library interface {
	AddBook(title, authorName string, price float64, ebookAvailable bool, publishDate time.Time) error
	ChangeName(idString, title string) (int, error)
	FindBook(idString string) (error, []byte)
	DeleteBook(idString string) error
	FindBooksByParams(title, authorName string, priceRangeValues []string) ([]Book, error)
	RetrieveStore() ([]interface{}, int64, error)
}
