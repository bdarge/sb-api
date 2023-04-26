package models

// Role Model. User has many roles
type Role struct {
	Model
	Name string `json:"name"`
}
