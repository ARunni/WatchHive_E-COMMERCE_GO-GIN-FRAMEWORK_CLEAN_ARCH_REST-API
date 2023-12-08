package interfaces

import (
	"WatchHive/pkg/domain"
	"WatchHive/pkg/utils/models"
	"mime/multipart"
)

type ProductUseCase interface {
	AddProduct(product models.AddProducts, file *multipart.FileHeader) (models.ProductResponse, error)
	ListProducts(int, int) ([]models.ProductUserResponse, error)
	EditProduct(product domain.Product) (models.ProductUserResponse, error)
	DeleteProduct(id string) error
	UpdateProduct(productID int, stock int) (models.ProductResponse, error)
}
