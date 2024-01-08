package responses

import (
	"github.com/gin-gonic/gin"
	"github.com/woonmapao/order-detail-service-go/models"
)

func CreateError(errors []string) gin.H {
	return gin.H{
		"status":  "error",
		"message": "validation failed",
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

func GetSuccess(data interface{}) gin.H {
	return gin.H{
		"status":  "success",
		"message": "get success",
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
