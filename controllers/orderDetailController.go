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

	// Extract updated order detail data from the request body
	var updatedData models.OrderDetail
	err := c.ShouldBindJSON(&updatedData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()},
		)
		return
	}

	// Validate the input data
	err = validators.ValidateUpdatedOrderDetailData(updatedData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update the order detail in the database
	var existingOrderDetail models.OrderDetail
	err = initializer.DB.First(&existingOrderDetail, orderDetailID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Order detail not found",
		})
		return
	}

	initializer.DB.Model(&existingOrderDetail).Updates(updatedData)

	// Return a JSON response with the updated order detail
	c.JSON(http.StatusOK, gin.H{
		"updatedOrderDetail": existingOrderDetail,
	})
}

func DeleteOrderDetail(c *gin.Context) {
	// Extract order detail ID from the request parameters
	orderDetailID := c.Param("id")

	// Delete the order detail from the database
	err := initializer.DB.Delete(&models.OrderDetail{}, orderDetailID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete order detail",
		})
		return
	}

	// Return a JSON response indicating success
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func GetOrderDetailsByOrderID(c *gin.Context) {
	// Extract order ID from the request parameters
	orderID := c.Param("id")

	// Query the database for order details associated with the order
	orderDetails, err := func(orderId string) ([]models.OrderDetail, error) {

		var orderDetails []models.OrderDetail
		err := initializer.DB.Where("order_id = ?", orderID).Find(&orderDetails).Error
		if err != nil {
			return nil, err
		}
		return orderDetails, nil
	}(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch order details",
		})
		return
	}

	// Return a JSON response with the order details
	c.JSON(http.StatusOK, gin.H{"orderDetails": orderDetails})
}
