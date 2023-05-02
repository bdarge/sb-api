package models

// Business Model.
type Business struct {
	Model
	Name       string  `json:"name"`
	HourlyRate uint32  `gorm:"column:hourly_rate" json:"hourlyRate"`
	Vat        uint32  `json:"vat"`
	Address    Address `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"address"`
	Users      []User  `json:"users"`
}
