package models

import (
	"time"
)

type Coupon struct {
	CouponName      string    `json:"coupon_name" gorm:"unique"`
	OfferPercentage int       `json:"offer_percentage" gorm:"not null"`
	ExpireDate      time.Time `json:"expire_date" gorm:"type:date"`
}

type CouponResp struct {
	ID uint `json:"id"`
	OfferPercentage int `json:"offer_percentage"`
	ExpireDate time.Time `json:"expire_date"`
	
}