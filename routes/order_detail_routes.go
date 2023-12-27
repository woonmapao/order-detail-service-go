package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/woonmapao/order-detail-service-go/controllers"
)

func SetupOrderDetailRoutes(router *gin.Engine) {
	orderDetailGroup := router.Group("/order-details")
	{
		orderDetailGroup.GET("/", controllers.GetAllOrderDetails)
		orderDetailGroup.GET("/:id", controllers.GetOrderDetailByID)
		orderDetailGroup.POST("/", controllers.AddOrderDetail)
		orderDetailGroup.PUT("/:id", controllers.UpdateOrderDetail)
		orderDetailGroup.DELETE("/:id", controllers.DeleteOrderDetail)
		orderDetailGroup.GET("/orders/:id", controllers.GetOrderDetailsByOrderID)
	}
}
