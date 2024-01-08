package controllers

import (
	"errors"

	m "github.com/woonmapao/order-detail-service-go/models"
	"gorm.io/gorm"
)

func GetDetails(db *gorm.DB) (*[]m.OrderDetail, error) {

	var details []m.OrderDetail
	err := db.Find(&details).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &details, errors.New("failed to fetch details")
	}
	if err == gorm.ErrRecordNotFound {
		return &details, errors.New("no details found")
	}
	return &details, nil
}
