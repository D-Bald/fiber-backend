package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User struct
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Username  string             `bson:"username" json:"username" `
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"password"`
	Names     string             `bson:"names,omitempty" json:"names" `
}
