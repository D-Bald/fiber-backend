package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Content   struct
type Conten  t struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt 	time.Time          `bson:"created_at"`
	UpdatedAt 	time.Time          `bson:"updated_at"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description,omitempty" json:"description"`
	Text        string             `bson:"text" json:"text"`
	Tags		[]string			`bson:"tags,omitempty" json:"tags"`
}
