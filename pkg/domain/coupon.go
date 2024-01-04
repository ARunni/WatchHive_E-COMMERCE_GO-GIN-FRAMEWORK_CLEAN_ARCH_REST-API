package domain

import (
	"time"

	"gorm.io/gorm"
)

type Coupon struct {
	gorm.Model
	ID              uint      `json:"id" gorm:"primaryKey"`
	CouponName      string    `json:"coupon_name" gorm:"unique"`
	OfferPercentage int       `json:"offer_percentage" gorm:"not null"`
	ExpireDate      time.Time `json:"expire_date" gorm:"type:date"`
}
