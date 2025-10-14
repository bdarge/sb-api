package models

// Lang represents supported Language
type Lang struct {
	ID 				string `gorm:"column:id" json:"id"`
	Language  string `gorm:"column:language" json:"language"`
	Currency  string `gorm:"column:currency" json:"currency"`
}


// Langs list of Lang
type Langs struct {
	Data  []Lang `json:"data"`
}
