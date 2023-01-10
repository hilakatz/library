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
