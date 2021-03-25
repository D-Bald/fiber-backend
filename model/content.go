package model

import (
	"github.com/jinzhu/gorm"
)

// Content struct
type Content struct {
	gorm.Model
	BaseModel
	Text string `gorm:"not null" json:"text"`
}
