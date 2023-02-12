package models

import (
	"time"
)

// Account for each user
type Account struct {
	ID        uint32    `json:"id,string"` // https://stackoverflow.com/a/21152548
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
	Email     *string   `gorm:"unique;not null" json:"email"`
	Password  *string   `gorm:"not null" json:"-"`
	UserID    int       `json:"userId"`
}
