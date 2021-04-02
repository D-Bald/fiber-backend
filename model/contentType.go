package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User struct
type ContentType struct {
	ID          primitive.ObjectID     `bson:"_id" json:"id" xml:"id" form:"email"`
	CreatedAt   time.Time              `bson:"created_at"`
	UpdatedAt   time.Time              `bson:"updated_at"`
	TypeName    string                 `bson:"typename" json:"typename" xml:"typename" form:"typename"`
	Collection  string                 `bson:"collection" json:"collection" xml:"collection" form:"collection"`
	FieldSchema map[string]interface{} `bson:"fields_schema" json:"field_schema" xml:"field_schema" form:"field_schema"` // not used yet. Could be used for introducing a Validator on Collection Creation.
}

func (ct *ContentType) Init() {
	ct.ID = primitive.NewObjectID()
	ct.CreatedAt = time.Now()
	ct.UpdatedAt = time.Now()
}
