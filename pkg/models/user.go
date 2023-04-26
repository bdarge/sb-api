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
	UserName      string        `json:"username"`
	HourlyRate    string        `gorm:"column:hourly_rate" json:"hourlyRate"`
	BusinessName  string        `gorm:"column:business_name" json:"businessName"`
	Street        string        `json:"street"`
	PostalCode    string        `gorm:"column:postal_code" json:"postalCode"`
	City          string        `json:"city"`
	Country       string        `json:"country"`
	LandLinePhone string        `gorm:"column:landline_phone" json:"landlinePhone"`
	MobilePhone   string        `gorm:"column:mobile_phone" json:"mobilePhone"`
	Vat           string        `json:"vat"`
	Transactions  []Transaction `gorm:"foreignKey:CreatedBy" json:"transactions"`
	Account       Account
	Roles         []Role `gorm:"many2many:user_roles;" json:"roles"`
}
