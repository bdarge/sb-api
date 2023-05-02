package models

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint32         `json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}

// User Model. User has an account
type User struct {
	Model
	UserName     string        `json:"username"`
	Address      Address       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"address"`
	Roles        []Role        `gorm:"many2many:user_roles;" json:"roles"`
	Transactions []Transaction `gorm:"foreignKey:CreatedBy" json:"transactions"`
	BusinessID   uint32        `json:"businessId"`
	Business     Business      `json:"business"`
	Account      Account
}
