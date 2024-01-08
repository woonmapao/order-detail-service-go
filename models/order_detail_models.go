package models

import (
	"gorm.io/gorm"
)

type OrderDetail struct {
	gorm.Model
	OrderID   int     `json:"order_id"`   // Foreign key to Orders
	ProductID int     `json:"product_id"` // Foreign key to Product
	Quantity  int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
}

type DetailRequest struct {
	OrderID   int `json:"order_id"`   // Foreign key to Orders
	ProductID int `json:"product_id"` // Foreign key to Product
	Quantity  int `json:"quantity"`
}
