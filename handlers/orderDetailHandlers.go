package handlers

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	ctrl "github.com/woonmapao/order-detail-service-go/controllers"
	"github.com/woonmapao/order-detail-service-go/initializer"
	i "github.com/woonmapao/order-detail-service-go/initializer"
	"github.com/woonmapao/order-detail-service-go/models"
	m "github.com/woonmapao/order-detail-service-go/models"
	"github.com/woonmapao/order-detail-service-go/responses"
	r "github.com/woonmapao/order-detail-service-go/responses"
	s "github.com/woonmapao/order-detail-service-go/services"
)

const productServiceURL = "http://localhost:2002/products"

func GetDetailsHandler(c *gin.Context) {

	details, err := ctrl.GetDetails(i.DB)
	if err != nil {
		c.JSON(http.StatusNotFound,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	c.JSON(http.StatusOK,
		r.GetDetailsSuccess(*details),
	)
}

func GetDetailHandler(c *gin.Context) {

	id, err := ctrl.GetID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	detail, err := ctrl.GetDetail(id, i.DB)
	if err != nil {
		c.JSON(http.StatusNotFound,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	c.JSON(http.StatusOK,
		responses.GetSuccess(detail),
	)
}

func AddDetailHandler(c *gin.Context) {

	var body m.DetailRequest
	err := ctrl.BindAndValidate(c, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Start a transaction
	tx, err := ctrl.StartTrx(c)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				tx.Error.Error(),
			}))
		return
	}

	var wg sync.WaitGroup
	var product *m.Product
	var orderErr, productErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		_, orderErr = s.GetOrder(body.OrderID)
	}()

	go func() {
		defer wg.Done()
		product, productErr = s.GetProduct(body.ProductID)
	}()

	wg.Wait()

	if orderErr != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				orderErr.Error(),
			}))
		return
	}
	if productErr != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				productErr.Error(),
			}))
		return
	}

	adding := m.OrderDetail{
		OrderID:   body.OrderID,
		ProductID: body.ProductID,
		Quantity:  body.Quantity,
		Subtotal:  product.Price * float64(body.Quantity),
	}
	err = ctrl.AddDetail(&adding, tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Commit the transaction and check for commit errors
	err = ctrl.CommitTrx(c, tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Return a JSON response with the newly created order detail
	c.JSON(http.StatusOK,
		responses.CreateSuccess(adding),
	)
}

func UpdateOrderDetailHandler(c *gin.Context) {

	id, err := ctrl.GetID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	var body m.DetailRequest
	err = ctrl.BindAndValidate(c, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Start a transaction
	tx, err := ctrl.StartTrx(c)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				tx.Error.Error(),
			}))
		return
	}

	// Validate if ProductID and OrderID exist
	var wg sync.WaitGroup
	var product *m.Product
	var orderErr, productErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		_, orderErr = s.GetOrder(body.OrderID)
	}()

	go func() {
		defer wg.Done()
		product, productErr = s.GetProduct(body.ProductID)
	}()

	wg.Wait()

	if orderErr != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				orderErr.Error(),
			}))
		return
	}
	if productErr != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				productErr.Error(),
			}))
		return
	}

	exist, err := ctrl.GetDetail(id, tx)
	if err != nil {
		c.JSON(http.StatusNotFound,
			responses.CreateError([]string{
				err.Error(),
			}))
		return
	}

	updating := m.OrderDetail{
		OrderID:   body.OrderID,
		ProductID: body.ProductID,
		Quantity:  body.Quantity,
		Subtotal:  product.Price * float64(body.Quantity),
	}
	err = ctrl.UpdateOrderDetail(&updating, exist, tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Commit the transaction and check for commit errors
	err = ctrl.CommitTrx(c, tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.UpdateSuccess(&updating),
	)

}

func DeleteDetailHandler(c *gin.Context) {

	id, err := ctrl.GetID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Start a transaction
	tx, err := ctrl.StartTrx(c)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				tx.Error.Error(),
			}))
		return
	}

	// Check if its exist
	_, err = ctrl.GetDetail(id, i.DB)
	if err != nil {
		c.JSON(http.StatusNotFound,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Delete the order detail from the database
	err = ctrl.DeleteDetail(id, tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Commit the transaction and check for commit errors
	err = ctrl.CommitTrx(c, tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			r.CreateError([]string{
				err.Error(),
			}))
		return
	}

	// Return a JSON response indicating success
	c.JSON(http.StatusOK, r.DeleteSuccess())
}

func GetOrderDetailsByOrderID(c *gin.Context) {
	// Extract order ID from the request parameters
	orderID := c.Param("id")

	id, err := strconv.Atoi(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid order ID",
				err.Error(),
			}))
		return
	}

	// Query the database for order details associated with the order
	var orderDetails []models.OrderDetail
	err = initializer.DB.Where("order_id = ?", id).Find(&orderDetails).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch order details",
				err.Error(),
			}))
		return
	}
	if len(orderDetails) == 0 {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"No order details found for the give order ID",
			}))
		return
	}

	// Return a JSON response with the order details
	c.JSON(http.StatusOK,
		responses.GetSuccessResponseForMultipleOrderDetails(orderDetails))
}
