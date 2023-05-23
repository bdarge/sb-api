package models

// Business Model.
type Business struct {
	Model
	Name       string `json:"name"`
	HourlyRate uint32 `gorm:"column:hourly_rate" json:"hourlyRate"`
	Vat        uint32 `json:"vat"`
	Users      []User `json:"users"`
	Street     string `json:"street"`
	PostalCode string `gorm:"column:postal_code" json:"postalCode"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	Landline   string `json:"landline"`
	Mobile     string `json:"mobile"`
}
