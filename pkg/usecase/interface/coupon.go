package interfaces

import "WatchHive/pkg/utils/models"

type CouponUsecase interface {
	AddCoupon(coupon models.Coupon)(models.CouponResp,error)
	GetCoupon()([]models.CouponResp,error)
	EditCoupon(coupon models.CouponResp) (models.CouponResp,error)
}