package models

import (
	"time"
)

// Customer Model
type Customer struct {
	ID           uint32    `json:"id,string"` // https://stackoverflow.com/a/21152548
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	DeletedAt    time.Time `json:"deletedAt"`
	Email        *string   `gorm:"unique;not null"`
	Name         *string   `gorm:"not null"`
	Dispositions []Disposition
}
