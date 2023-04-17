package models

// Account for each user
type Account struct {
	Model
	Email    *string `gorm:"unique;not null" json:"email"`
	Password *string `gorm:"not null" json:"-"`
	UserID   uint32  `json:"userId"`
}
