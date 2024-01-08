package responses

import (
	"github.com/gin-gonic/gin"
	"github.com/woonmapao/order-detail-service-go/models"
)

func CreateErrorResponse(errors []string) gin.H {
	return gin.H{
		"status":  "error",
		"message": "Validation failed",
		"data": gin.H{
			"error": errors,
		},
	}
}

func CreateSuccessResponse(data interface{}) gin.H {
	return gin.H{
		"status":  "success",
		"message": "Request successful",
		"data":    data,
	}
}

func GetSuccessResponse(data interface{}) gin.H {
	return gin.H{
		"status":  "success",
		"message": "Request successful",
		"data":    data,
	}
}

func GetDetailsSuccess(orderDetails []models.OrderDetail) gin.H {
	return gin.H{
		"status":  "success",
		"message": "order details retrieved successfully",
		"data": gin.H{
			"orderDetails": orderDetails,
		},
	}
}
