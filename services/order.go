package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	m "github.com/woonmapao/order-detail-service-go/models"
)

const orderServiceURL = "http://localhost:3030/order"

type OrderResponse struct {
	Data    OrderData `json:"data"`
	Message string    `json:"message"`
	Status  string    `json:"status"`
}

type OrderData struct {
	Order m.Order `json:"order"`
}

func GetOrder(id int) (*m.Order, error) {

	// Make get request
	resp, err := http.Get(fmt.Sprintf("%s/%d", orderServiceURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the status code of the response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"failed to fetch order. Status code: %d",
			resp.StatusCode,
		)
	}

	// Decode the JSON response into a ProductResponse struct
	var orderResponse OrderResponse
	err = json.NewDecoder(resp.Body).Decode(&orderResponse)
	if err != nil {
		return nil, err
	}

	// Return the product from the nested structure
	return &orderResponse.Data.Order, nil
}
