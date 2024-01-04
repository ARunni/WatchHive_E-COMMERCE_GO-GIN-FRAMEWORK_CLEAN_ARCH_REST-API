package usecase

import (
	helper "WatchHive/pkg/helper/interface"
	repo "WatchHive/pkg/repository/interface"
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/errmsg"
	"WatchHive/pkg/utils/models"
	"errors"
)

type couponUsecase struct {
	couponRepo repo.CouponRepository
	h          helper.Helper
}

func NewCouponUsecase(couponRep repo.CouponRepository, h helper.Helper) interfaces.CouponUsecase {
	return &couponUsecase{
		couponRepo: couponRep,
		h:          h,
	}
}

func (cu *couponUsecase) AddCoupon(coupon models.Coupon) (models.CouponResp, error) {
	if coupon.CouponName == "" {
		return models.CouponResp{}, errors.New(errmsg.ErrFieldEmpty)
	}
	if coupon.OfferPercentage <= 0 {
		return models.CouponResp{}, errors.New(errmsg.ErrDataZero)
	}
	formattedExpireDate := coupon.ExpireDate.Format("02-01-2006")
	ok := cu.h.ValidateDate(formattedExpireDate)
	if !ok {
		return models.CouponResp{}, errors.New(errmsg.ErrInvalidDate)
	}
	couponResp, err := cu.couponRepo.AddCoupon(coupon)
	if err != nil {
		return models.CouponResp{}, err
	}
	return couponResp, nil
}
