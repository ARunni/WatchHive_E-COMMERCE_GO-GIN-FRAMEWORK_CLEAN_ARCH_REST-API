run :
	go run ./cmd/api

wire: ## Generate wire_gen.go
	cd pkg/di && wire
swag :
	swag init -g cmd/api/main.go -o ./cmd/api/docs