package domain

import "time"

type ProductOffer struct {
	ID                 uint      `json:"id" gorm:"unique; not null"`
	ProductID          uint      `json:"product_id"`
	Products           Product   `json:"-" gorm:"foreignkey:ProductID"`
	OfferName          string    `json:"offer_name"`
	DiscountPercentage int       `json:"discount_percentage"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
}
type CategoryOffer struct {
	ID                 uint      `json:"id" gorm:"unique; not null"`
	CategoryID         uint      `json:"category_id"`
	Category           Category  `json:"-" gorm:"foreignkey:CategoryID"`
	OfferName          string    `json:"offer_name"`
	DiscountPercentage int       `json:"discount_percentage"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
}
