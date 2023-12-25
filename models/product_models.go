package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name          string  `json:"name"`
	Category      string  `json:"category"`
	Price         float64 `json:"price"`
	Description   string  `json:"description"`
	StockQuantity int     `json:"stockQuantity"`
	ReorderLevel  int     `json:"reorderLevel"`
}
