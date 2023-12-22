// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"WatchHive/pkg/api"
	"WatchHive/pkg/api/handler"
	"WatchHive/pkg/config"
	"WatchHive/pkg/db"
	"WatchHive/pkg/helper"
	"WatchHive/pkg/repository"
	"WatchHive/pkg/usecase"
)

// Injectors from wire.go:

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	gormDB, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}
	adminRepository := repository.NewAdminRepository(gormDB)
	interfacesHelper := helper.NewHelper(cfg)
	adminUseCase := usecase.NewAdminUseCase(adminRepository, interfacesHelper)
	adminHandler := handler.NewAdminHandler(adminUseCase)
	userRepository := repository.NewUserRepository(gormDB)
	walletRepository := repository.NewWalletRepository(gormDB)
	userUseCase := usecase.NewUserUseCase(userRepository, cfg, interfacesHelper, walletRepository)
	userHandler := handler.NewUserHandler(userUseCase)
	otpRepository := repository.NewOtpRepository(gormDB)
	otpUseCase := usecase.NewOtpUseCase(cfg, otpRepository, interfacesHelper)
	otpHandler := handler.NewOtpHandler(otpUseCase)
	categoryRepository := repository.NewCategoryRepository(gormDB)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)
	productRepository := repository.NewProductRepository(gormDB)
	productUseCase := usecase.NewProductUseCase(productRepository, interfacesHelper, categoryRepository)
	productHandler := handler.NewProductHandler(productUseCase)
	cartRepository := repository.NewCartRepository(gormDB)
	cartUseCase := usecase.NewCartUseCase(cartRepository, productRepository)
	cartHandler := handler.NewCartHandler(cartUseCase)
	orderRepository := repository.NewOrderRepository(gormDB)
	paymentRepository := repository.NewPaymentRepository(gormDB)
	orderUseCase := usecase.NewOrderUseCase(orderRepository, walletRepository, cartRepository, userRepository, paymentRepository)
	paymentUseCase := usecase.NewPaymentUseCase(paymentRepository, orderRepository, cfg)
	orderHandler := handler.NewOrderHandler(orderUseCase, paymentUseCase)
	paymentHandler := handler.NewPaymentHandler(paymentUseCase)
	walletUsecase := usecase.NewWalletUsecase(walletRepository)
	walletHandler := handler.NewWalletHandler(walletUsecase)
	serverHTTP := http.NewServerHTTP(adminHandler, userHandler, otpHandler, categoryHandler, productHandler, cartHandler, orderHandler, paymentHandler, walletHandler)
	return serverHTTP, nil
}
