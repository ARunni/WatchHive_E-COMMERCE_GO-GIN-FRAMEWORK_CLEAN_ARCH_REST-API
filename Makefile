run :
	go run ./cmd/api

wire: ## Generate wire_gen.go
	cd pkg/di && wire