package main

import (
	"log"

	"github.com/woonmapao/order-detail-service-go/initializer"
	"github.com/woonmapao/order-detail-service-go/models"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.DBInitializer()
}

func main() {

	err := initializer.DB.AutoMigrate(&models.OrderDetail{})
	if err != nil {
		log.Fatal("Failed to perform auto migration: &v", err)
	}
}
