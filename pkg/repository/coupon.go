package repository

import (
	interfaces "WatchHive/pkg/repository/interface"
	"WatchHive/pkg/utils/errmsg"
	"WatchHive/pkg/utils/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type couponRepo struct {
	DB *gorm.DB
}

func NewCouponRepository(db *gorm.DB) interfaces.CouponRepository {
	return &couponRepo{DB: db}
}

func (cr *couponRepo) AddCoupon(coupon models.Coupon) (models.CouponResp, error) {
	var CouponResp models.CouponResp
	dateNow := time.Now()
	querry := `
	insert into coupons
	 (coupon_name, offer_percentage, expire_date,created_at)
	  values(?,?,?,?) returning *

	`
	result := cr.DB.Raw(querry, coupon.CouponName, coupon.OfferPercentage, coupon.ExpireDate, dateNow).Scan(&CouponResp)
	if result.Error != nil {
		return models.CouponResp{}, errors.New(errmsg.ErrWriteDB)
	}

	return CouponResp, nil

}
func (cr *couponRepo) IsCouponExistByName(couponName string) (bool, error) {
	var count int
	err := cr.DB.Raw("select count(*) from coupons where coupon_name = ?", couponName).Scan(&count).Error
	if err != nil {
		return false, errors.New(errmsg.ErrGetDB)
	}
	return count > 0, nil
}
func (cr *couponRepo) IsCouponExistByID(couponID int) (bool, error) {
	var count int
	err := cr.DB.Raw("select count(*) from coupons where id = ?", couponID).Scan(&count).Error
	if err != nil {
		return false, errors.New(errmsg.ErrGetDB)
	}
	return count > 0, nil
}

func (cr *couponRepo) GetCoupon() ([]models.CouponResp, error) {
	var coupon []models.CouponResp
	err := cr.DB.Raw("select * from coupons").Scan(&coupon).Error
	if err != nil {
		return []models.CouponResp{}, errors.New(errmsg.ErrGetDB)
	}
	return coupon, nil
}

func (cr *couponRepo) GetCouponData(couponID int) (models.CouponResp, error) {
	var coupon models.CouponResp
	err := cr.DB.Raw("select * from coupons").Scan(&coupon).Error
	if err != nil {
		return models.CouponResp{}, errors.New(errmsg.ErrGetDB)
	}
	return coupon, nil
}

func (cr *couponRepo) EditCoupon(coupon models.CouponResp) (models.CouponResp, error) {
	curretTime := time.Now()
	querrry := `
	update coupons set
	 updated_at = ?, coupon_name = ?,
	  offer_percentage = ?, expire_date = ? 
	  where id = ?
	`
	err := cr.DB.Exec(querrry, curretTime, coupon.CouponName, coupon.OfferPercentage, coupon.ExpireDate, coupon.ID).Error
	if err != nil {
		return models.CouponResp{}, errors.New(errmsg.ErrWriteDB)
	}
	return coupon, nil
}

func (cr *couponRepo) GetCouponIdByName(couponName string) (int, error) {
	var id int
	err := cr.DB.Raw("select id from coupons where coupon_name = ?", couponName).Scan(&id).Error
	if err != nil {
		return 0, errors.New(errmsg.ErrGetDB)
	}
	return id, nil
}
