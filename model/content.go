package model

import (
	"time"
)

// Content   struct
type Content struct {
	ID            string                 `bson:"_id" json:"_id" xml:"_id" form:"_id" query:"_id"`
	CreatedAt     time.Time              `bson:"created_at" json:"created_at" xml:"created_at" form:"created_at" query:"created_at"`
	UpdatedAt     time.Time              `bson:"updated_at" json:"updated_at" xml:"updated_at" form:"updated_at" query:"updated_at"`
	ContentTypeID string                 `bson:"content_type_id" json:"content_type_id" xml:"content_type_id" form:"content_type_id"`
	Title         string                 `bson:"title" json:"title" xml:"title" form:"title" query:"title"`
	Published     *bool                  `bson:"published" json:"published" xml:"published" form:"published" query:"published"`
	Tags          []string               `bson:"tags" json:"tags" xml:"tags" form:"tags" query:"tags"`
	Fields        map[string]interface{} `bson:"fields,inline" json:"fields" xml:"fields" form:"fields" query:"fields"`
}

// Initialize metadata (ID is initialized by database controller)
func (c *Content) Init(ct ContentType) {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	c.ContentTypeID = ct.ID
}

// Fields that can be updated through API endpoints
type ContentUpdate struct {
	Title     string                 `bson:"title,omitempty" json:"title" xml:"title" form:"title"`
	Published *bool                  `bson:"published,omitempty" json:"published" xml:"published" form:"published"` // empty value is `nil` pointer, so it can be differentiated from `false` value for `omitempty`flag
	Tags      []string               `bson:"tags,omitempty" json:"tags" xml:"tags" form:"tags"`
	Fields    map[string]interface{} `bson:"fields,inline,omitempty" json:"fields" xml:"fields" form:"fields"`
}
