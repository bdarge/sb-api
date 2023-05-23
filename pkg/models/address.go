package models

// Address Model
type Address struct {
	Model
	Street     string `json:"street"`
	PostalCode string `gorm:"column:postal_code" json:"postalCode"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	Landline   string `json:"landline"`
	Mobile     string `json:"mobile"`
	UserID     uint32 `json:"userId"`
}

// CustomerAddress Model
type CustomerAddress struct {
	Model
	Street     string `json:"street"`
	PostalCode string `gorm:"column:postal_code" json:"postalCode"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	Landline   string `json:"landline"`
	Mobile     string `json:"mobile"`
	CustomerID uint32 `json:"customerId"`
}
