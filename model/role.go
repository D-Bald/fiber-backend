package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Saves the user-IDs of users with this role
type Role struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id" xml:"_id" form:"_id"`
	Role   string             `bson:"role" json:"role" xml:"role" form:"role"`
	Weight uint               `bson:"weight" json:"weight" xml:"weight" form:"weight"`
}

// Initialize metadata
func (r *Role) Init() {
	r.ID = primitive.NewObjectID()
}
