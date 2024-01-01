package interfaces

import (
	"WatchHive/pkg/domain"
	"WatchHive/pkg/utils/models"
)

type OfferRepository interface {
	AddProductOffer(ProductOffer models.ProductOfferResp) error
	GetProductOffer() ([]domain.ProductOffer, error)
	AddCategoryOffer(CategoryOffer models.CategorytOfferResp) error
	GetCategoryOffer() ([]domain.CategoryOffer, error)
	ExpireCategoryOffer(id int) error
	ExpireProductOffer(id int) error
	GetCatOfferPercent(categoryId int) (int,error)
	GetProOfferPercent(productId int) (int,error)
}
