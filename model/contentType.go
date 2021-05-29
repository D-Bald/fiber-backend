package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ContentType struct
type ContentType struct {
	ID          primitive.ObjectID              `bson:"_id" json:"_id" xml:"_id" form:"_id"`
	CreatedAt   time.Time                       `bson:"created_at"`
	UpdatedAt   time.Time                       `bson:"updated_at"`
	TypeName    string                          `bson:"typename" json:"typename" xml:"typename" form:"typename"`
	Collection  string                          `bson:"collection" json:"collection" xml:"collection" form:"collection"`
	Permissions map[string][]primitive.ObjectID `bson:"permissions" json:"permissions" xml:"permissions" form:"permissions"`
	FieldSchema map[string]interface{}          `bson:"field_schema" json:"field_schema" xml:"field_schema" form:"field_schema"` // not used yet. Could be used for introducing a Validator on Collection Creation.
}

// Initialize metadata
func (ct *ContentType) Init() {
	ct.ID = primitive.NewObjectID()
	ct.CreatedAt = time.Now()
	ct.UpdatedAt = time.Now()
}

// Fields that can be updated through API endpoints
type ContentTypeUpdate struct {
	TypeName    string                 `bson:"typename,omitempty" json:"typename" xml:"typename" form:"typename"`
	Collection  string                 `bson:"collection,omitempty" json:"collection" xml:"collection" form:"collection"`
	Permissions map[string][]string    `bson:"permissions,omitempty" json:"permissions" xml:"permissions" form:"permissions"`
	FieldSchema map[string]interface{} `bson:"field_schema,omitempty" json:"field_schema" xml:"field_schema" form:"field_schema"`
}
