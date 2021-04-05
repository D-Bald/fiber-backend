package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User struct
type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id" xml:"id" form:"id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Username  string             `bson:"username" json:"username" xml:"username" form:"username"`
	Email     string             `bson:"email" json:"email" xml:"email" form:"email"`
	Password  string             `bson:"password" json:"password" xml:"password" form:"password"`
	Names     string             `bson:"names" json:"names" xml:"names" form:"names"`
	Roles     []string           `bson:"roles" json:"roles" xml:"roles" form:"roles"`
}

// Initialize metadata
func (u *User) Init() {
	u.ID = primitive.NewObjectID()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.Roles = []string{"user"}
}

// Fields that can be updated through API endpoints
type UpdateUserInput struct {
	Username string   `bson:"username,omitempty" json:"username" xml:"username" form:"username"`
	Email    string   `bson:"email,omitempty" json:"email" xml:"email" form:"email"`
	Password string   `bson:"password,omitempty" json:"password" xml:"password" form:"password"`
	Names    string   `bson:"names,omitempty" json:"names" xml:"names" form:"names"`
	Roles    []string `bson:"roles,omitempty" json:"roles" xml:"roles" form:"roles"`
}
