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
		OrderID   int `json:"orderId"`
		ProductID int `json:"productId"`
		Quantity  int `json:"quantity" binding:"gte=1"`
	}
	err = c.ShouldBindJSON(&updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid request format",
				err.Error(),
			}))
		return
	}

	// Check for empty values
	if updateData.OrderID == 0 || updateData.ProductID == 0 || updateData.Quantity == 0 {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Username, email, and password are required fields",
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
	}

	// Get existing order detail from the database
	var orderDetail models.OrderDetail
	err = tx.First(&orderDetail, id).Error
	if err != nil {
		tx.Callback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch order detail",
				err.Error(),
			}))
		return
	}
	if orderDetail == (models.OrderDetail{}) {
		tx.Rollback()
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"Order detail not found",
			}))
	}

	// Validate the update data
	err = validations.ValidateOrderDetailData(
		updateData.OrderID, updateData.ProductID,
		updateData.Quantity, tx)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				err.Error(),
			}))
		return
	}

	// Update order detail fields based on the provided update data
	// Only update the fields that are present in the request
	if updateData.OrderID != 0 {
		orderDetail.OrderID = updateData.OrderID
	}
	if updateData.ProductID != 0 {
		orderDetail.ProductID = updateData.ProductID
	}
	if updateData.Quantity != 0 {
		orderDetail.Quantity = updateData.Quantity
	}
	if updateData.Quantity != 0 || updateData.ProductID != 0 {
		// Fetch the product price from the product-service
		productPrice, err := getProductPrice(updateData.ProductID)
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
		subtotal := float64(updateData.Quantity) * productPrice
		orderDetail.Subtotal = subtotal
	}

	// Save the updated order detail to the database
	err = tx.Save(&orderDetail).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to update order detail",
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
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponse(&orderDetail),
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
				err.Error(),
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

	// Check if the order detail with the give ID exists
	var orderDetail models.OrderDetail
	err = tx.First(&orderDetail, id).Error
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

	err = tx.Delete(&orderDetail, id).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to delete order detail",
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
				err.Error(), // Include the specific error message
			}))
		return
	}

	// Return a JSON response indicating success
	c.JSON(http.StatusOK, responses.CreateSuccessResponse(&orderDetail))
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
