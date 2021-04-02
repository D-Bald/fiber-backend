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
	Role      string             `bson:"role" json:"role" xml:"role" form:"role"`
}

func (u *User) Init() {
	u.ID = primitive.NewObjectID()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.Role = "user"
}
