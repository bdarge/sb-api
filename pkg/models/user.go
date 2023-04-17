package models

// User Model. User has an account
type User struct {
	Model
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
	Transactions  []Transaction `gorm:"foreignKey:CreatedBy" json:"transactions"`
	Account       Account
}
