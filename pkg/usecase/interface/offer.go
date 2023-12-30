package interfaces

import (
	"WatchHive/pkg/domain"
	"WatchHive/pkg/utils/models"
)

type OfferUsecase interface {
	AddProductOffer(ProductOffer models.ProductOfferResp) error
	GetProductOffer() ([]domain.ProductOffer, error)
	ExpireProductOffer(id int) error
	AddCategoryOffer(CategoryOffer models.CategorytOfferResp) error
	GetCategoryOffer() ([]domain.CategoryOffer, error)
	ExpireCategoryOffer(id int) error
}
