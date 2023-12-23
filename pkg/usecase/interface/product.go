package interfaces

import (
	"WatchHive/pkg/utils/models"
	"mime/multipart"
)

type ProductUseCase interface {
	AddProduct(product models.AddProducts, file *multipart.FileHeader) (models.ProductResponse, error)
	ListProducts(int, int) ([]models.ProductUserResponse, error)
	EditProduct(product models.ProductEdit) (models.ProductUserResponse, error)
	DeleteProduct(id string) error
	UpdateProduct(productID int, stock int) (models.ProductResponse, error)
}
