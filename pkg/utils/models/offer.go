package models

type ProductOfferResp struct {
	ProductID          uint   `json:"product_id" binding:"required"`
	OfferName          string `json:"offer_name" binding:"required"`
	DiscountPercentage int    `json:"discount_percentage" binding:"required"`
}

type CategorytOfferResp struct {
	CategoryID         uint   `json:"category_id" binding:"required"`
	OfferName          string `json:"offer_name" binding:"required"`
	DiscountPercentage int    `json:"discount_percentage" binding:"required"`
}
