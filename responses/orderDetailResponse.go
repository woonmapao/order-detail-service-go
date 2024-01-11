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

func CreateSuccess(data interface{}) gin.H {
	return gin.H{
		"status":  "success",
		"message": "create success",
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

func UpdateSuccess(data interface{}) gin.H {
	return gin.H{
		"status":  "success",
		"message": "update success",
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
