package validations

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ValidateOrderDetailData(data struct {
	OrderID   int     `json:"orderId" binding:"required"`
	ProductID int     `json:"productId" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,gte=1"`
	Subtotal  float64 `json:"subtotal"`
}) error {

	// Fetch information using HTTP call
	orderProductInfo, err := getOrderAndProductInfo(data.OrderID, data.ProductID)
	if err != nil {
		return fmt.Errorf("failed to retrieve order and product information: %v", err)
	}

	// Check if the order exists
	if orderProductInfo.OrderStatusCode != http.StatusOK {
		return fmt.Errorf("order with ID %d does not exist", data.OrderID)
	}

	// Check if the product exists
	if orderProductInfo.ProductStatusCode != http.StatusOK {
		return fmt.Errorf("product with ID %d does not exist", data.ProductID)
	}

	// Check if the product has sufficient stock
	if orderProductInfo.StockQuantity < data.Quantity {
		return fmt.Errorf("insufficient stock for product with ID %d", data.ProductID)
	}

	return nil

}

func ValidateOrderDetailUpdateData(updateData struct {
	Quantity int     `json:"quantity" binding:"required,gte=1"`
	Subtotal float64 `json:"subtotal"`
}) error {

	if updateData.Quantity <= 0 {
		return fmt.Errorf("quantity must be a positive integer")
	}

	if updateData.Subtotal < 0 {
		return fmt.Errorf("subtotal must be a non-negative value")
	}

	return nil
}

// Make an HTTP request to both the order and product services
func getOrderAndProductInfo(orderID, productID int) (*OrderProductInfo, error) {
	// Make HTTP GET request to both order and product services concurrently
	// My first time trying go routines in a project :D
	orderCh := make(chan *http.Response)
	productCh := make(chan *http.Response)

	go func() {
		resp, err := http.Get(fmt.Sprintf("http://order-service/api/orders/%d", orderID))
		if err != nil {
			orderCh <- nil
		} else {
			orderCh <- resp
		}
	}()

	go func() {
		resp, err := http.Get(fmt.Sprintf("http://product-service/api/products/%d", productID))
		if err != nil {
			productCh <- nil
		} else {
			productCh <- resp
		}
	}()

	// Wait for responses from both services
	orderResp := <-orderCh
	productResp := <-productCh

	// Close the channels
	close(orderCh)
	close(productCh)

	// Parse the responses and return the combined information
	_, productInfo, err := parseOrderAndProductInfo(orderResp, productResp)
	if err != nil {
		return nil, err
	}

	return &OrderProductInfo{
		OrderStatusCode:   orderResp.StatusCode,
		ProductStatusCode: productResp.StatusCode,
		StockQuantity:     productInfo.StockQuantity,
	}, nil
}

// parseOrderAndProductInfo parses the responses from the order and product services.
func parseOrderAndProductInfo(orderResp, productResp *http.Response) (*OrderInfo, *ProductInfo, error) {
	// Parse order information
	var orderInfo OrderInfo
	if orderResp != nil {
		err := json.NewDecoder(orderResp.Body).Decode(&orderInfo)
		if err != nil {
			return nil, nil, err
		}
		defer orderResp.Body.Close()
	}

	// Parse product information
	var productInfo ProductInfo
	if productResp != nil {
		err := json.NewDecoder(productResp.Body).Decode(&productInfo)
		if err != nil {
			return nil, nil, err
		}
		defer productResp.Body.Close()
	}

	return &orderInfo, &productInfo, nil
}

// OrderInfo represents information about an order.
type OrderInfo struct {
	// No use right now
}

// ProductInfo represents information about a product.
type ProductInfo struct {
	StockQuantity int `json:"stockQuantity"`
	// Only need this for now
}

// OrderProductInfo combines information from the order and product services.
type OrderProductInfo struct {
	OrderStatusCode   int
	ProductStatusCode int
	StockQuantity     int
}
