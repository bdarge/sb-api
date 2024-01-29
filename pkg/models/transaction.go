package models

import (
	"time"
)

// Transaction Model
type Transaction struct {
	Model
	Description   string            `gorm:"column:description" json:"description"`
	DeliveryDate  *time.Time        `gorm:"column:delivery_date" json:"deliveryDate"`
	InvoiceNumber string            `gorm:"column:invoice_number" json:"invoiceNumber"`
	Currency      string            `gorm:"column:currency" json:"currency"`
	Items         []TransactionItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
	CreatedBy     uint32            `json:"createdBy"`
	CustomerID    uint32            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"customerId"`
	Customer      Customer          `json:"customer"`
	RequestType   string            `gorm:"column:request_type" json:"requestType"`
}

// TransactionItem Model
type TransactionItem struct {
	Model
	Description   string  `gorm:"column:description" json:"description"`
	Qty           uint32  `gorm:"column:qty" json:"qty"`
	Unit          string  `gorm:"column:unit" json:"unit"`
	UnitPrice     float64 `gorm:"column:unit_price" json:"unitPrice"`
	TransactionID uint32  `gorm:"column:transaction_id" json:"transactionId"`
}

// Transactions Model
type Transactions struct {
	Limit uint32        `json:"limit"`
	Page  uint32        `json:"page"`
	Total uint32        `json:"total"`
	Data  []Transaction `json:"data"`
}

// TransactionItems Model
type TransactionItems struct {
	Limit uint32            `json:"limit"`
	Page  uint32            `json:"page"`
	Total uint32            `json:"total"`
	Data  []TransactionItem `json:"data"`
}
