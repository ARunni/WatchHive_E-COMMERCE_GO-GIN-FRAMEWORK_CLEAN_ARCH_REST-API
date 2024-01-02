package usecase

import (
	"WatchHive/pkg/domain"
	interfaces "WatchHive/pkg/repository/interface"
	services "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/errmsg"
	"WatchHive/pkg/utils/models"
	"errors"
)

type offerUsecase struct {
	repo interfaces.OfferRepository
}

func NewOferUsecase(repo interfaces.OfferRepository) services.OfferUsecase {
	return &offerUsecase{
		repo: repo,
	}
}

func (ou *offerUsecase) AddProductOffer(ProductOffer models.ProductOfferResp) error {
	if ProductOffer.ProductID <= 0 {
		return errors.New(errmsg.ErrDataZero)
	}
	if ProductOffer.DiscountPercentage <= 0 {
		return errors.New(errmsg.ErrDataZero)
	}
	if ProductOffer.OfferName == ""{
		return errors.New(errmsg.ErrFieldEmpty)
	}

	if err := ou.repo.AddProductOffer(ProductOffer); err != nil {
		return errors.New(errmsg.ErrOfferAdd)
	}
	return nil
}
func (ou *offerUsecase) AddCategoryOffer(CategoryOffer models.CategorytOfferResp) error {
	if CategoryOffer.CategoryID <= 0 {
		return errors.New(errmsg.ErrDataZero)
	}
	if CategoryOffer.DiscountPercentage <= 0 {
		return errors.New(errmsg.ErrDataZero)
	}
	if CategoryOffer.OfferName == ""{
		return errors.New(errmsg.ErrFieldEmpty)
	}
	if err := ou.repo.AddCategoryOffer(CategoryOffer); err != nil {
		return errors.New(errmsg.ErrOfferAdd)
	}
	return nil
}
func (ou *offerUsecase) GetProductOffer() ([]domain.ProductOffer, error) {
	offers, err := ou.repo.GetProductOffer()
	if err != nil {
		return []domain.ProductOffer{}, err
	}
	return offers, nil
}

func (ou *offerUsecase) GetCategoryOffer() ([]domain.CategoryOffer, error) {
	offers, err := ou.repo.GetCategoryOffer()
	if err != nil {
		return []domain.CategoryOffer{}, err
	}
	return offers, nil
}

func (ou *offerUsecase) ExpireProductOffer(id int) error {
	ok, err := ou.repo.IsProductOfferExist(id)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(errmsg.ErrNotExist)
	}
	if err := ou.repo.ExpireProductOffer(id); err != nil {
		return err
	}

	return nil
}

func (ou *offerUsecase) ExpireCategoryOffer(id int) error {
	ok, err := ou.repo.IsCategoryOfferExist(id)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(errmsg.ErrNotExist)
	}
	if err := ou.repo.ExpireCategoryOffer(id); err != nil {
		return err
	}

	return nil
}
