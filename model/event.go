package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Event struct
type Event struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description,omitempty" json:"description"`
	Date        string             `bson:"date" json:"date"` // time.time muss man parsen.
}
