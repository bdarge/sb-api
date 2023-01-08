package models

import "gorm.io/gorm"

// Disposition Model
type Disposition struct {
	gorm.Model
	Description       string            `gorm:"column:description" json:"description"`
	DeliveryDate      string            `gorm:"column:delivery_date" json:"deliveryDate"`
	InvoiceNumber     string            `gorm:"column:invoice_number" json:"invoiceNumber"`
	Currency          string            `gorm:"column:currency" json:"currency"`
	DispositionNumber string            `gorm:"column:disposition_number" json:"dispositionNumber"`
	Items             []DispositionItem `json:"items"`
	CreatedBy         uint              `json:"createdBy"`
	CustomerID        uint              `json:"customerId"`
	RequestType       string            `gorm:"column:requestType" json:"requestType"`
}

// DispositionItem Model
type DispositionItem struct {
	gorm.Model
	Description   string `gorm:"column:description" json:"description"`
	Qty           string `gorm:"column:qty" json:"qty"`
	Unit          string `gorm:"column:unit" json:"unit"`
	UnitPrice     string `gorm:"column:unit_price" json:"unitPrice"`
	DispositionID uint   `json:"dispositionId"`
}

type Page struct {
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	Total int64 `json:"total"`
}

type Dispositions struct {
	Page
	Data []Disposition `json:"data"`
}
