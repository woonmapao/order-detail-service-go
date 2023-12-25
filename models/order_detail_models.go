package models

import (
	"gorm.io/gorm"
)

type OrderDetail struct {
	gorm.Model
	OrderID   int     `json:"orderId"`   // Foreign key to Orders
	ProductID int     `json:"productId"` // Foreign key to Product
	Quantity  int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
}
