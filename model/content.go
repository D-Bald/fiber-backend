package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Content   struct
type Content struct {
	ID          primitive.ObjectID     `bson:"_id" json:"id" xml:"id" form:"id"`
	CreatedAt   time.Time              `bson:"created_at"`
	UpdatedAt   time.Time              `bson:"updated_at"`
	ContentType primitive.ObjectID     `bson:"content_type" json:"content_type" xml:"content_type" form:"content_type"`
	Title       string                 `bson:"title" json:"title" xml:"title" form:"title"`
	Published   bool                   `bson:"published" json:"published" xml:"published" form:"published"`
	Tags        []string               `bson:"tags" json:"tags" xml:"tags" form:"tags"`
	Fields      map[string]interface{} `bson:"fields,inline" json:"fields" xml:"fields" form:"fields"` // inline Flag used, to handle fields of a bson.D input as if they were part of Content Document (https://pkg.go.dev/go.mongodb.org/mongo-driver/bson#hdr-Structs)
}
