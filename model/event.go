package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Event struct
type Event struct {
	gorm.Model
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	Date        time.Time `gorm:"not null" json:"date"`
}
