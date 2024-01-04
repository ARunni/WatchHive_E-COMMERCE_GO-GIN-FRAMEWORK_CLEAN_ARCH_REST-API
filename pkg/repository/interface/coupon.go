package interfaces

import "WatchHive/pkg/utils/models"

type CouponRepository interface {
	AddCoupon(models.Coupon) (models.CouponResp, error)
	IsCouponExistByName(couponName string) (bool, error)
	IsCouponExistByID(couponID int) (bool, error)
}
