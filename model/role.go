package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Saves the user-IDs of users with this role
type Role struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id" xml:"_id" form:"_id"`
	Tag  string             `bson:"tag,omitempty" json:"tag" xml:"tag" form:"tag"`
	Name string             `bson:"name,omitempty" json:"name" xml:"name" form:"name"`
}

// Initialize metadata
func (r *Role) Init() {
	r.ID = primitive.NewObjectID()
}
