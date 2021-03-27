package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User struct
type ContentType struct {
	ID         primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	CreatedAt  time.Time              `bson:"created_at"`
	UpdatedAt  time.Time              `bson:"updated_at"`
	TypeName   string                 `bson:"type_name"`
	Collection string                 `bson:"collection"`
	published  bool                   `bson:"published"`
	Fields     map[string]interface{} `bson:"fields"`
}
