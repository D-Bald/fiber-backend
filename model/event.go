package model

import (
	"github.com/jinzhu/gorm"
)

// Event struct
type Event struct {
	gorm.Model
	Title       string `gorm:"not null" json:"title"`
	Description string `json:"description"`
	Date        string `gorm:"not null" json:"date"` // time.time muss man parsen.
}
