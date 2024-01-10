package interfaces

import "WatchHive/pkg/utils/models"

type CouponRepository interface {
	AddCoupon(models.Coupon) (models.CouponResp, error)
	IsCouponExistByName(couponName string) (bool, error)
	GetCouponIdByName(couponName string) (int, error)
	IsCouponExistByID(couponID int) (bool, error)
	GetCoupon()([]models.CouponResp,error)
	GetCouponData(couponID int)(models.CouponResp,error)
	EditCoupon(coupon models.CouponResp)(models.CouponResp,error)
}
