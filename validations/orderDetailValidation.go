package validations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/woonmapao/order-detail-service-go/models"
	"gorm.io/gorm"
)

const (
	orderServiceURL   = "http://localhost:3030/order"
	productServiceURL = "http://localhost:2002/products"
)

func ValidateOrderDetailData(orderID, productID, quantity int, tx *gorm.DB) error {

	// Fetch information using HTTP call
	orderProductInfo, err := getOrderAndProductInfo(orderID, productID)
	if err != nil {
		return fmt.Errorf(
			"failed to retrieve order and product information: %v", err,
		)
	}

	// Check if the order exists
	if orderProductInfo.OrderID != 0 {
		return fmt.Errorf(
			"order with ID %d does not exist", orderID,
		)
	}

	// Check if the product exists
	if orderProductInfo.ProductID != 0 {
		return fmt.Errorf(
			"product with ID %d does not exist", productID,
		)
	}

	return nil

}

// Make an HTTP request to both the order and product services
func getOrderAndProductInfo(orderID, productID int) (*OrderProduct, error) {
	// Make HTTP GET request to both order and product services concurrently
	// My first time trying go routines in a project :D
	orderCh := make(chan *http.Response)
	productCh := make(chan *http.Response)
	errCh := make(chan error, 2)

	// HTTP GET concurrently
	go func() {
		resp, err := http.Get(fmt.Sprintf("%s/%d", orderServiceURL, orderID))
		if err != nil {
			errCh <- err
		} else {
			orderCh <- resp
		}
	}()

	go func() {
		resp, err := http.Get(fmt.Sprintf("%s/%d", productServiceURL, productID))
		if err != nil {
			errCh <- err
		} else {
			productCh <- resp
		}
	}()

	// Wait for responses from both services or an error
	var orderResp, productResp *http.Response
	for i := 0; i < 2; i++ {
		select {
		case orderResp = <-orderCh:
		case productResp = <-productCh:
		case err := <-errCh:
			close(orderCh)
			close(productCh)
			close(errCh)
			return nil, fmt.Errorf("failed to make HTTP request: %v", err)
		}
	}
	close(orderCh)
	close(productCh)
	close(errCh)

	// Parse the responses and return the combined information
	order, product, err := parseOrderAndProductInfo(orderResp, productResp)
	if err != nil {
		return nil, err
	}

	return &OrderProduct{
		OrderID:              int(order.Model.ID),
		ProductID:            int(product.Model.ID),
		ProductStockQuantity: product.StockQuantity,
	}, nil
}

func parseOrderAndProductInfo(orderResp, productResp *http.Response) (*models.Order, *models.Product, error) {

	var order models.Order
	var product models.Product

	// Parse order information
	if orderResp != nil {
		defer orderResp.Body.Close()

		if orderResp.StatusCode == http.StatusOK {
			// Decode the JSON response into the orderInfo struct
			err := json.NewDecoder(orderResp.Body).Decode(&order)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to decode order response: %v", err)
			}
		} else {
			// Handle non-OK status code (if needed)
			return nil, nil, fmt.Errorf("order service returned non-OK status: %d", orderResp.StatusCode)
		}
	}

	// Parse the product response
	if productResp != nil {
		defer productResp.Body.Close()

		if productResp.StatusCode == http.StatusOK {
			// Decode the JSON response into the productInfo struct
			err := json.NewDecoder(productResp.Body).Decode(&product)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to decode product response: %v", err)
			}
		} else {
			// Handle non-OK status code (if needed)
			return nil, nil, fmt.Errorf("product service returned non-OK status: %d", productResp.StatusCode)
		}
	}

	return &order, &product, nil
}

type OrderProduct struct {
	// Order fields
	OrderID     int
	UserID      int       `json:"userId"` // Foreign key to User
	OrderDate   time.Time `json:"orderDate"`
	TotalAmount float64   `json:"totalAmount"`
	Status      string    `json:"status"`

	// Product fields
	ProductID            int
	ProductName          string  `json:"productName"`
	ProductCategory      string  `json:"productCategory"`
	ProductPrice         float64 `json:"productPrice"`
	ProductDescription   string  `json:"productDescription"`
	ProductStockQuantity int     `json:"productStockQuantity"`
	ProductReorderLevel  int     `json:"productReorderLevel"`
}
