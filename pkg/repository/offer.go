package repository

import (
	"WatchHive/pkg/domain"
	interfaces "WatchHive/pkg/repository/interface"
	"WatchHive/pkg/utils/errmsg"
	"WatchHive/pkg/utils/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type OfferRepository struct {
	DB *gorm.DB
}

func NewOfferRepository(DB *gorm.DB) interfaces.OfferRepository {
	return &OfferRepository{DB: DB}
}
func (or *OfferRepository) AddProductOffer(ProductOffer models.ProductOfferResp) error {
	var count int
	err := or.DB.Raw("select count(*) from product_offers where offer_name =? and product_id=?", ProductOffer.OfferName, ProductOffer.ProductID).Scan(&count).Error
	if err != nil {
		return errors.New("error in getting offer details")
	}
	if count > 0 {
		return errors.New("offer already exist")
	}

	err = or.DB.Raw("SELECT COUNT(*) FROM product_offers WHERE product_id = ?", ProductOffer.ProductID).Scan(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		err = or.DB.Exec("DELETE FROM product_offers WHERE product_id = ?", ProductOffer.ProductID).Error
		if err != nil {
			return err
		}
	}

	startDate := time.Now()
	endDate := time.Now().Add(time.Hour * 24 * 5)
	err = or.DB.Exec("INSERT INTO product_offers (product_id, offer_name, discount_percentage, start_date, end_date) VALUES (?, ?, ?, ?, ?)", ProductOffer.ProductID, ProductOffer.OfferName, ProductOffer.DiscountPercentage, startDate, endDate).Error
	if err != nil {
		return err
	}

	return nil
}

func (or *OfferRepository) AddCategoryOffer(CategoryOffer models.CategorytOfferResp) error {
	var count int
	err := or.DB.Raw("select count(*) from category_offers where offer_name =? and category_id=?", CategoryOffer.OfferName, CategoryOffer.CategoryID).Scan(&count).Error
	if err != nil {
		return errors.New("error in getting offer details")
	}
	if count > 0 {
		return errors.New("offer already exist")
	}

	err = or.DB.Raw("SELECT COUNT(*) FROM category_offers WHERE category_id = ?", CategoryOffer.CategoryID).Scan(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		err = or.DB.Exec("DELETE FROM category_offers WHERE category_id = ?", CategoryOffer.CategoryID).Error
		if err != nil {
			return err
		}
	}

	startDate := time.Now()
	endDate := time.Now().Add(time.Hour * 24 * 5)
	err = or.DB.Exec("INSERT INTO category_offers (category_id, offer_name, discount_percentage, start_date, end_date) VALUES (?, ?, ?, ?, ?)", CategoryOffer.CategoryID, CategoryOffer.OfferName, CategoryOffer.DiscountPercentage, startDate, endDate).Error
	if err != nil {
		return err
	}

	return nil
}

func (or *OfferRepository) GetProductOffer() ([]domain.ProductOffer, error) {
	var productOfferDetails []domain.ProductOffer
	err := or.DB.Raw("SELECT * FROM product_offers").Scan(&productOfferDetails).Error
	if err != nil {
		return []domain.ProductOffer{}, errors.New("error in getting product offers")
	}
	return productOfferDetails, nil
}

func (or *OfferRepository) GetCategoryOffer() ([]domain.CategoryOffer, error) {
	var categoryOfferDetails []domain.CategoryOffer
	err := or.DB.Raw("SELECT * FROM category_offers").Scan(&categoryOfferDetails).Error
	if err != nil {
		return []domain.CategoryOffer{}, errors.New("error in getting category offers")
	}
	return categoryOfferDetails, nil
}
func (of *OfferRepository) ExpireProductOffer(id int) error {
	if err := of.DB.Exec("DELETE FROM product_offers WHERE id = $1", id).Error; err != nil {
		return err
	}

	return nil
}
func (of *OfferRepository) ExpireCategoryOffer(id int) error {
	if err := of.DB.Exec("DELETE FROM category_offers WHERE id = $1", id).Error; err != nil {
		return err
	}

	return nil
}
 func (or *OfferRepository) GetCatOfferPercent(categoryId int) (int,error) {
	var percent int
	err := or.DB.Raw("select discount_percentage from category_offers where category_id = ? ",categoryId).Scan(&percent).Error
	if err != nil {
		return 0,errors.New(errmsg.ErrGetDB)
	}
	return percent,nil
 }
 func (or *OfferRepository) GetProOfferPercent(productId int) (int,error) {
	var percent int
	err := or.DB.Raw("select discount_percentage from category_offers where product_id = ? ",productId).Scan(&percent).Error
	if err != nil {
		return 0,errors.New(errmsg.ErrGetDB)
	}
	return percent,nil
 }