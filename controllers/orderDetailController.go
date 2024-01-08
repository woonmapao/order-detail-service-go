package controllers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
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

func GetID(c *gin.Context) (int, error) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, errors.New("invalid user id")
	}
	return id, nil
}

func GetDetail(id int, db *gorm.DB) (*m.OrderDetail, error) {

	var detail m.OrderDetail
	err := db.First(&detail, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &detail, errors.New("something went wrong")
	}
	if err == gorm.ErrRecordNotFound {
		return &detail, errors.New("order detail not found")
	}
	return &detail, nil
}
