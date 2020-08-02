package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type OneTimeScript struct {
	gorm.Model
	Token  string    `gorm:"not null"`
	Expire time.Time `gorm:"not null"`
	Email  string    `gorm:"not null"`
}
