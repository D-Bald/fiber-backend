package model

import (
	"time"
)

// ContentType struct
type ContentType struct {
	ID          string                 `bson:"_id" json:"_id" xml:"_id" form:"_id"`
	CreatedAt   time.Time              `bson:"created_at"`
	UpdatedAt   time.Time              `bson:"updated_at"`
	TypeName    string                 `bson:"typename" json:"typename" xml:"typename" form:"typename"`
	Collection  string                 `bson:"collection" json:"collection" xml:"collection" form:"collection"`
	Permissions map[string][]Role      `bson:"permissions" json:"permissions" xml:"permissions" form:"permissions"`     // key: method that will be allowed; key: list of role tags
	FieldSchema map[string]interface{} `bson:"field_schema" json:"field_schema" xml:"field_schema" form:"field_schema"` // not used yet. Could be used for introducing a Validator on Collection Creation.
}

// Initialize metadata (ID is initialized by database controller)
func (ct *ContentType) Init() {
	ct.CreatedAt = time.Now()
	ct.UpdatedAt = time.Now()
}

// Fields that can be updated through API endpoints
type ContentTypeInput struct {
	TypeName    string                 `bson:"typename,omitempty" json:"typename" xml:"typename" form:"typename"`
	Collection  string                 `bson:"collection,omitempty" json:"collection" xml:"collection" form:"collection"`
	Permissions map[string][]Role      `bson:"permissions,omitempty" json:"permissions" xml:"permissions" form:"permissions"`
	FieldSchema map[string]interface{} `bson:"field_schema,omitempty" json:"field_schema" xml:"field_schema" form:"field_schema"`
}
