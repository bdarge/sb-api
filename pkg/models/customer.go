package models

import "gorm.io/gorm"

// Customer Model
type Customer struct {
	gorm.Model
	Email        *string `gorm:"unique;not null"`
	Name         *string `gorm:"not null"`
	Dispositions []Disposition
}
