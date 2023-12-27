package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/woonmapao/order-detail-service-go/models"
)

const productServiceURL = "http://localhost:2002/products"

type ProductResponse struct {
	Data    ProductData `json:"data"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
}

type ProductData struct {
	Product models.Product `json:"product"`
}

func GetProductByID(productID int) (*models.Product, error) {
	// Make a GET request to the product-service API
	resp, err := http.Get(fmt.Sprintf("%s/%d", productServiceURL, productID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the status code of the response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch product. Status code: %d", resp.StatusCode)
	}

	// Decode the JSON response into a ProductResponse struct
	var productResponse ProductResponse
	err = json.NewDecoder(resp.Body).Decode(&productResponse)
	if err != nil {
		return nil, err
	}

	// Return the product from the nested structure
	return &productResponse.Data.Product, nil
}
