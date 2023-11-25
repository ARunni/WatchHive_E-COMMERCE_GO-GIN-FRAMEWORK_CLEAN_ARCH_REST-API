package interfaces

import (
	"WatchHive/pkg/domain"
	"WatchHive/pkg/utils/models"
)

type ProductUseCase interface {
	AddProduct(inventory models.AddProducts) (models.ProductResponse, error)
	ListProducts(int, int) ([]models.ProductUserResponse, error)
	EditProduct(domain.Product, int) (domain.Product, error)
	DeleteProduct(id string) error
	UpdateProduct(productID int, stock int) (models.ProductResponse, error)
}
