package models

// Customer Model
type Customer struct {
	Model
	Email   *string         `json:"email" gorm:"unique;not null"`
	Name    *string         `json:"name" gorm:"not null"`
	Address CustomerAddress `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"address"`
}

type Customers struct {
	Limit uint32     `json:"limit"`
	Page  uint32     `json:"page"`
	Total uint32     `json:"total"`
	Data  []Customer `json:"data"`
}
