package models

// Address Model
type Address struct {
	Model
	Street        string `json:"street"`
	PostalCode    string `gorm:"column:postal_code" json:"postalCode"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	LandLinePhone string `gorm:"column:landline_phone" json:"landlinePhone"`
	MobilePhone   string `gorm:"column:mobile_phone" json:"mobilePhone"`
	UserID        uint32 `json:"userId"`
	BusinessID    uint32 `json:"businessId"`
}
