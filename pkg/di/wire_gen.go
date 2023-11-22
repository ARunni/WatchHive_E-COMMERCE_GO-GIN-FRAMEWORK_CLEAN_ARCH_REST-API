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
	userUseCase := usecase.NewUserUseCase(userRepository, cfg, interfacesHelper)
	userHandler := handler.NewUserHandler(userUseCase)
	otpRepository := repository.NewOtpRepository(gormDB)
	otpUseCase := usecase.NewOtpUseCase(cfg, otpRepository, interfacesHelper)
	otpHandler := handler.NewOtpHandler(otpUseCase)
	categoryRepository := repository.NewCategoryRepository(gormDB)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)
	serverHTTP := http.NewServerHTTP(adminHandler, userHandler, otpHandler, categoryHandler)
	return serverHTTP, nil
}
