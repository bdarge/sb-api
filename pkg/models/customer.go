package models

import (
	"gorm.io/gorm"
	"time"
)

// Customer Model
type Customer struct {
	ID        uint32         `json:"id,string"` // https://stackoverflow.com/a/21152548
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
	Email     *string        `json:"email" gorm:"unique;not null"`
	Name      *string        `json:"name" gorm:"not null"`
}

type Customers struct {
	Limit uint32     `json:"limit"`
	Page  uint32     `json:"page"`
	Total uint32     `json:"total"`
	Data  []Customer `json:"data"`
}
