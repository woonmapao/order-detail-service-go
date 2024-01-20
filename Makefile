build:
	@go build -o bin/order-detail-service

run: build
	@./bin/order-detail-service

test:
	 @go test -v ./...