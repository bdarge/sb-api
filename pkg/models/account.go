package models

import (
	"gorm.io/gorm"
)

// Account for each user
type Account struct {
	gorm.Model
	Email    *string `gorm:"unique;not null" json:"email"`
	Password *string `gorm:"not null" json:"-"`
	UserID   int     `json:"userId"`
}
