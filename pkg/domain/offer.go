package domain

import "time"

type ProductOffer struct {
	ID                 uint      `json:"id" gorm:"unique; not null"`
	ProductID          uint      `json:"product_id" gorm:"not null"`
	Products           Product   `json:"-" gorm:"foreignkey:ProductID"`
	OfferName          string    `json:"offer_name" gorm:"not null"`
	DiscountPercentage int       `json:"discount_percentage" gorm:"not null"`
	StartDate          time.Time `json:"start_date" gorm:"not null"`
	EndDate            time.Time `json:"end_date" gorm:"not null"`
}
type CategoryOffer struct {
	ID                 uint      `json:"id" gorm:"unique; not null"`
	CategoryID         uint      `json:"category_id" gorm:"not null"`
	Category           Category  `json:"-" gorm:"foreignkey:CategoryID"`
	OfferName          string    `json:"offer_name" gorm:"not null"`
	DiscountPercentage int       `json:"discount_percentage" gorm:"not null"`
	StartDate          time.Time `json:"start_date" gorm:"not null"`
	EndDate            time.Time `json:"end_date" gorm:"not null"`
}
