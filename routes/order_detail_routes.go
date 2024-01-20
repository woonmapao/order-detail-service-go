package routes

import (
	"github.com/gin-gonic/gin"
	h "github.com/woonmapao/order-detail-service-go/handlers"
)

func SetupOrderDetailRoutes(router *gin.Engine) {

	orderDetailGroup := router.Group("/order-details")
	{
		orderDetailGroup.POST("/", h.AddDetailHandler)

		orderDetailGroup.GET("/", h.GetDetailsHandler)
		orderDetailGroup.GET("/:id", h.GetDetailHandler)

		orderDetailGroup.PUT("/:id", h.UpdateOrderDetailHandler)

		orderDetailGroup.DELETE("/:id", h.DeleteDetailHandler)
	}
}
