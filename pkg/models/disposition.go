package models

import (
	"time"
)

type Model struct {
	ID        uint32    `json:"id,string"` // https://stackoverflow.com/a/21152548
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

// Disposition Model
type Disposition struct {
	Model
	Description   string            `gorm:"column:description" json:"description"`
	DeliveryDate  time.Time         `gorm:"column:delivery_date" json:"deliveryDate"`
	InvoiceNumber string            `gorm:"column:invoice_number" json:"invoiceNumber"`
	Currency      string            `gorm:"column:currency" json:"currency"`
	Items         []DispositionItem `json:"items"`
	CreatedBy     uint32            `json:"createdBy"`
	CustomerID    uint32            `json:"customerId"`
	RequestType   string            `gorm:"column:requestType" json:"requestType"`
}

// DispositionItem Model
type DispositionItem struct {
	Model
	Description   string  `gorm:"column:description" json:"description"`
	Qty           uint32  `gorm:"column:qty" json:"qty"`
	Unit          string  `gorm:"column:unit" json:"unit"`
	UnitPrice     float64 `gorm:"column:unit_price" json:"unitPrice"`
	DispositionID uint32  `json:"dispositionId"`
}

type Dispositions struct {
	Limit uint32        `json:"limit"`
	Page  uint32        `json:"page"`
	Total uint32        `json:"total"`
	Data  []Disposition `json:"data"`
}
