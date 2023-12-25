package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID      int       `json:"userId"` // Foreign key to User
	OrderDate   time.Time `json:"orderDate"`
	TotalAmount float64   `json:"totalAmount"`
	Status      string    `json:"status"`
}
