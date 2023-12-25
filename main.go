package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/woonmapao/order-detail-service-go/initializer"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.DBInitializer()
}

func main() {

	r := gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)

}
