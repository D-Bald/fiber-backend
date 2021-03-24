package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Content struct
type Content struct {
	gorm.Model
	Type        string    `gorm:"not null" json:"type"`
	Title       string    `gorm:"not null" json:"title"`
	Date        time.Time `gorm:"not null" json:"date"`
	Description string    `json:"description"`
	Text        string    `gorm:"not null" json:"amount"`
	Author      User      `gorm:"embedded"`
}
