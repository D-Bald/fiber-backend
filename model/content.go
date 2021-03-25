package model

import (
	"github.com/jinzhu/gorm"
)

// Content struct
type Content struct {
	gorm.Model
	Title       string `gorm:"not null" json:"title"`
	Description string `json:"description"`
	Text        string `gorm:"not null" json:"text"`
}
