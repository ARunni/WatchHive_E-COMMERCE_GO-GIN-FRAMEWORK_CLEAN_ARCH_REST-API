package interfaces

import (
	"WatchHive/pkg/domain"
	"WatchHive/pkg/utils/models"
)

type ProductRepository interface {
	AddProduct(inventory models.AddProducts, url string) (models.ProductResponse, error)
	ListProducts(int, int) ([]models.ProductUserResponse, error)
	EditProduct(domain.Product, int) (domain.Product, error)
	DeleteProduct(id string) error
	CheckProduct(pid int) (bool, error)
	UpdateProduct(pid int, stock int) (models.ProductResponse, error)
}
