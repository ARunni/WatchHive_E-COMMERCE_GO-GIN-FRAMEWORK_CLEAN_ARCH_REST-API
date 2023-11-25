//go:build wireinject
// +build wireinject

package di

import (
	http "WatchHive/pkg/api"
	"WatchHive/pkg/api/handler"
	"WatchHive/pkg/config"
	"WatchHive/pkg/db"
	"WatchHive/pkg/helper"
	"WatchHive/pkg/repository"
	"WatchHive/pkg/usecase"

	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(

		db.ConnectDatabase,

		repository.NewUserRepository,
		repository.NewAdminRepository,
		repository.NewOtpRepository,
		repository.NewCategoryRepository,
		repository.NewProductRepository,

		usecase.NewUserUseCase,
		usecase.NewAdminUseCase,
		usecase.NewOtpUseCase,
		usecase.NewCategoryUseCase,
		usecase.NewProductUseCase,

		handler.NewUserHandler,
		handler.NewAdminHandler,
		handler.NewOtpHandler,
		handler.NewCategoryHandler,
		handler.NewProductHandler,

		helper.NewHelper,

		http.NewServerHTTP,
	)

	return &http.ServerHTTP{}, nil
}
