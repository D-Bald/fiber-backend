package model

import (
	"github.com/jinzhu/gorm"
)

// BaseModel struct
type BaseModel struct {
	gorm.Model
	Title       string `gorm:"not null" json:"title"`
	Description string `json:"description"`
	Author      User   `gorm:"embedded"`
}
