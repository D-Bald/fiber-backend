package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Event struct
type Event struct {
	gorm.Model
	BaseModel
	Date time.Time `gorm:"not null" json:"date"`
}
