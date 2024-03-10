run :
	go run ./cmd/api

wire: ## Generate wire_gen.go
	cd pkg/di && wire
swag :
	swag init -g cmd/api/main.go -o ./cmd/api/docs

mock : ## mockgen
	mockgen -source pkg/repository/interface/user.go -destination pkg/repository/mock/user_mock.go -package mock

	mockgen -source pkg/usecase/interface/user.go -destination pkg/usecase/mock/user_mock.go -package mock
	mockgen -source pkg/repository/interface/admin.go -destination pkg/repository/mock/admin_mock.go -package mock
	mockgen -source pkg/helper/interface/helper.go -destination pkg/helper/mock/helper_mock.go -package mock
	mockgen -source pkg/repository/interface/wallet.go -destination pkg/repository/mock/wallet_mock.go -package mock

	
	