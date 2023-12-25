package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/woonmapao/order-detail-service-go/initializer"
	"github.com/woonmapao/order-detail-service-go/models"
	"github.com/woonmapao/order-detail-service-go/responses"
	"github.com/woonmapao/order-detail-service-go/validations"
)

func GetAllOrderDetails(c *gin.Context) {
	// Retrieve order details from the database
	var orderDetails []models.OrderDetail
	err := initializer.DB.Find(&orderDetails).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch order details",
			}))
	}

	// Check if no order details were found
	if len(orderDetails) == 0 {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"No order details found",
			}))
	}

	// Return a JSON response with the list of order details
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponseForMultipleOrderDetails(
			orderDetails,
		))
}

func GetOrderDetailByID(c *gin.Context) {
	// Extract order detail ID from the request parameters
	orderDetailID := c.Param("id")

	// Convert order detail ID to integer (validation)
	id, err := strconv.Atoi(orderDetailID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid order detail ID",
			}))
		return
	}

	// Get the order detail from the database
	var orderDetail models.OrderDetail
	err = initializer.DB.First(&orderDetail, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch order detail",
			}))
		return
	}

	// Check if the order detail was not found
	if orderDetail == (models.OrderDetail{}) {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"Order detail not found",
			}))
		return
	}

	// Return success response with order detail
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponse(&orderDetail))
}

func CreateOrderDetail(c *gin.Context) {
	// Extract data from the request body
	var body struct {
		OrderID   int     `json:"orderId" binding:"required"`
		ProductID int     `json:"productId" binding:"required"`
		Quantity  int     `json:"quantity" binding:"required,gte=1"`
		Subtotal  float64 `json:"subtotal"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				err.Error(),
			}))
		return
	}

	// Validate the input data
	err = validations.ValidateOrderDetailData(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create order detail in the database
	orderDetail := models.OrderDetail{
		OrderID:   body.OrderID,
		ProductID: body.ProductID,
		Quantity:  body.Quantity,
		Subtotal:  body.Subtotal,
	}

	err = initializer.DB.Create(&orderDetail).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to create order detail",
			}))
		return
	}

	// Return a JSON response with the newly created order detail
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponse(&orderDetail),
	)
}

func UpdateOrderDetail(c *gin.Context) {
	// Extract order detail ID from the request parameters
	orderDetailID := c.Param("id")

	// Convert order detail ID to integer (validations)
	id, err := strconv.Atoi(orderDetailID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid order detail ID",
			}))
		return
	}

	// Get existing order detail from the database
	var existingOrderDetail models.OrderDetail
	err = initializer.DB.First(&existingOrderDetail, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch order detail",
			}))
		return
	}

	// Get data from the request body
	var updateData struct {
		Quantity int     `json:"quantity" binding:"required,gte=1"`
		Subtotal float64 `json:"subtotal"`
	}
	err = c.ShouldBindJSON(&updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				err.Error(),
			}))
		return
	}

	// Validate the update data
	err = validations.ValidateOrderDetailUpdateData(updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				err.Error(),
			}))
		return
	}

	// Update order detail fields
	existingOrderDetail.Quantity = updateData.Quantity
	existingOrderDetail.Subtotal = updateData.Subtotal

	// Save the updated order detail to the database
	err = initializer.DB.Save(&existingOrderDetail).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to update order detail",
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponse(&existingOrderDetail),
	)

}

func DeleteOrderDetail(c *gin.Context) {
	// Extract order detail ID from the request parameters
	orderDetailID := c.Param("id")

	// Convert order detail ID to integer (validations)
	id, err := strconv.Atoi(orderDetailID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid order detail ID",
			}))
		return
	}

	// Check if the order detail with the give ID exists
	var orderDetail models.OrderDetail
	err = initializer.DB.First(&orderDetail, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch order detail",
			}))
		return
	}

	if orderDetail == (models.OrderDetail{}) {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"Order detail not found",
			}))
		return
	}

	// Delete the order detail from the database

	err = initializer.DB.Delete(&orderDetail, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to delete order detail",
			}))
		return
	}

	// Return a JSON response indicating success
	c.JSON(http.StatusOK, responses.CreateSuccessResponse(nil))
}

func GetOrderDetailsByOrderID(c *gin.Context) {
	// Extract order ID from the request parameters
	orderID := c.Param("id")

	id, err := strconv.Atoi(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid order ID",
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
		responses.CreateSuccessResponseForMultipleOrderDetails(orderDetails))
}
