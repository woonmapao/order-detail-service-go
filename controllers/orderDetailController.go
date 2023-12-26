package controllers

import (
	"encoding/json"
	"fmt"
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
				err.Error(),
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
		responses.GetSuccessResponseForMultipleOrderDetails(
			orderDetails,
		),
	)
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
				err.Error(),
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
				err.Error(),
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
		responses.GetSuccessResponse(&orderDetail),
	)
}

func AddOrderDetail(c *gin.Context) {
	// Extract data from the request body
	var body struct {
		OrderID   int `json:"orderId" binding:"required"`
		ProductID int `json:"productId" binding:"required"`
		Quantity  int `json:"quantity" binding:"required,gte=1"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid request format",
				err.Error(),
			}))
		return
	}

	// Check for empty values
	if body.OrderID == 0 || body.ProductID == 0 || body.Quantity == 0 {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"OrderID, ProductID, and Quantity are required fields",
			}))
		return
	}

	// Start a transaction
	tx := initializer.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to begin transaction",
				tx.Error.Error(),
			}))
		return
	}

	// Validate the input data
	err = validations.ValidateOrderDetailData(
		body.OrderID, body.ProductID, body.Quantity, tx,
	)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Fetch the product price from the product-service
	productPrice, err := getProductPrice(body.ProductID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch product price",
				err.Error(),
			}))
		return
	}

	// Calculate the subtotal
	subtotal := float64(body.Quantity) * productPrice

	// Create order detail in the database
	orderDetail := models.OrderDetail{
		OrderID:   body.OrderID,
		ProductID: body.ProductID,
		Quantity:  body.Quantity,
		Subtotal:  subtotal,
	}
	err = tx.Create(&orderDetail).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to create order detail",
				err.Error(),
			}))
		return
	}

	// Commit the transaction and check for commit errors
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to commit transaction",
				err.Error(),
			}))
		return
	}

	// Return a JSON response with the newly created order detail
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponse(&orderDetail),
	)
}

// Function to fetch product price from product-service
func getProductPrice(productID int) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("http://product-service/api/products/%d", productID))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf(
			"failed to fetch product price. Status code: %d", resp.StatusCode,
		)
	}

	var product models.Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return 0, err
	}

	return product.Price, nil
}

func UpdateOrderDetail(c *gin.Context) {

	orderDetailID := c.Param("id")

	// Convert order detail ID to integer (validations)
	id, err := strconv.Atoi(orderDetailID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid order detail ID",
				err.Error(),
			}))
		return
	}

	// Get data from the request body
	var updateData struct {
		OrderID   int `json:"orderId" binding:"required"`
		ProductID int `json:"productId" binding:"required"`
		Quantity  int `json:"quantity" binding:"required,gte=1"`
	}
	err = c.ShouldBindJSON(&updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				err.Error(),
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
