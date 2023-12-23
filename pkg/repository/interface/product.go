package interfaces

import (
	"WatchHive/pkg/utils/models"
)

type ProductRepository interface {
	AddProduct(inventory models.AddProducts, url string) (models.ProductResponse, error)
	ListProducts(int, int) ([]models.ProductUserResponse, error)
	EditProduct(product models.ProductEdit) (models.ProductUserResponse, error)
	DeleteProduct(id string) error
	CheckProduct(pid int) (bool, error)
	CheckProductAndCat(prdt string, cat int) bool
	UpdateProduct(pid int, stock int) (models.ProductResponse, error)
	CheckProductAvailable(product_id int) (bool, error)
	GetPriceOfProduct(product_id int) (float64, error)
}
