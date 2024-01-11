package controllers

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	i "github.com/woonmapao/order-detail-service-go/initializer"
	m "github.com/woonmapao/order-detail-service-go/models"
	r "github.com/woonmapao/order-detail-service-go/responses"
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

func AddDetail(detail *m.OrderDetail, tx *gorm.DB) error {

	adding := m.OrderDetail{
		OrderID:   detail.OrderID,
		ProductID: detail.ProductID,
		Quantity:  detail.Quantity,
		Subtotal:  detail.Subtotal,
	}
	err := tx.Create(&adding).Error
	if err != nil {
		return errors.New("failed to add order detail")
	}
	return nil
}

func BindAndValidate[T any](c *gin.Context, body *T) error {

	err := c.ShouldBindJSON(&body)
	if err != nil {
		return errors.New(
			"invalid request format",
		)
	}
	if reflect.ValueOf(body).IsNil() {
		return errors.New("missing fields")
	}

	return nil
}

func StartTrx(c *gin.Context) (*gorm.DB, error) {

	tx := i.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func CommitTrx(c *gin.Context, tx *gorm.DB) error {

	err := tx.Commit().Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				"Failed to commit transaction",
				err.Error(),
			}))
		return err
	}
	return nil
}

func UpdateOrderDetail(update *m.OrderDetail, exist *m.OrderDetail, tx *gorm.DB) error {

	exist.OrderID = update.OrderID
	exist.ProductID = update.ProductID
	exist.Quantity = update.Quantity
	exist.Subtotal = update.Subtotal

	err := tx.Save(&exist).Error
	if err != nil {
		return errors.New("failed to update order detail")
	}
	return nil
}

func DeleteDetail(id int, tx *gorm.DB) error {

	err := tx.Delete(&m.OrderDetail{}, id)
	if err != nil {
		return errors.New("failed to delete order detail")
	}
	return nil
}
