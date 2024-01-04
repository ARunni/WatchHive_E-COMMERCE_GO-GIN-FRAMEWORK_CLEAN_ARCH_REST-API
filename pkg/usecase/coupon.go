package usecase

import (
	helper "WatchHive/pkg/helper/interface"
	repo "WatchHive/pkg/repository/interface"
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/errmsg"
	"WatchHive/pkg/utils/models"
	"errors"
	"time"
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
	parsedStartDate, err := time.Parse("02-01-2006", coupon.ExpireDate)
	if err != nil {
		err := errors.New(errmsg.ErrFormat + " :expire_date")
		return models.CouponResp{}, err
	}

	isValid := !parsedStartDate.IsZero()
	if !isValid {
		err := errors.New(errmsg.ErrFormat + " :expire_date")
		return models.CouponResp{}, err
	}
	ok,err := cu.couponRepo.IsCouponExistByName(coupon.CouponName)
	if err != nil {
		return models.CouponResp{},err
	}
	if ok {
		return models.CouponResp{},errors.New(errmsg.ErrCouponExistTrue)
	}
	couponResp, err := cu.couponRepo.AddCoupon(coupon)
	if err != nil {
		return models.CouponResp{}, err
	}
	return couponResp, nil
}
