package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title          string             `json:"title" bson:"title" binding:"required"`
	AuthorName     string             `json:"author_name" bson:"author_name" binding:"required"`
	Price          float64            `json:"price" bson:"price" binding:"required"`
	EbookAvailable bool               `json:"ebook_available" bson:"ebook_available"`
	PublishDate    time.Time          `json:"publish_date" bson:"publish_date" binding:"required"`
}

type Library interface {
	AddBook(title, authorName string, price float64, ebookAvailable bool, publishDate time.Time) (string, error)
	ChangeName(idString, title string) (int, error)
	FindBook(idString string) (Book, error)
	DeleteBook(idString string) error
	FindBooksByParams(title, authorName string, priceMin, priceMax float64) ([]Book, error)
	RetrieveStore() (int, int, error)
}
