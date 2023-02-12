package models

import (
	"time"
)

// User Model
type User struct {
	ID            uint32        `json:"id,string"` // https://stackoverflow.com/a/21152548
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	DeletedAt     time.Time     `json:"deletedAt"`
	UserName      string        `json:"username"`
	HourlyRate    string        `json:"hourlyRate"`
	BusinessName  string        `json:"businessName"`
	Street        string        `json:"street"`
	PostalCode    string        `json:"postalCode"`
	City          string        `json:"city"`
	Country       string        `json:"country"`
	LandLinePhone string        `gorm:"column:landline_phone" json:"landlinePhone"`
	MobilePhone   string        `gorm:"column:mobile_phone" json:"mobilePhone"`
	Vat           string        `json:"vat"`
	Dispositions  []Disposition `gorm:"foreignKey:CreatedBy" json:"dispositions"`
	Account       Account
}
