package models

type CheckoutDetails struct {
	AddressInfoResponse []AddressInfoResponse
	Payment_Method      []PaymentDetails
	Cart                []Cart
	Total_Price         float64
}

type OrderFromCart struct {
	PaymentID uint `json:"payment_id" binding:"required"`
	AddressID uint `json:"address_id" binding:"required"`
	CouponID  uint `json:"coupon_id" `
}
type OrderSuccessResponse struct {
	OrderID        uint    `json:"order_id"`
	ShipmentStatus string  `json:"shipment_status"`
	Total          float64 `json:"total"`
	FinalPrice     float64 `json:"finalprice"`
}

type OrderDetails struct {
	OrderId        int
	FinalPrice     float64
	ShipmentStatus string
	PaymentStatus  string
}

type OrderIncoming struct {
	UserID    int `json:"user_id"`
	PaymentID int `json:"payment_id"`
	AddressID int `json:"address_id"`
	CouponID  int `json:"coupon_id"`
}

type OrderProductDetails struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
}

type FullOrderDetails struct {
	OrderDetails        OrderDetails
	OrderProductDetails []OrderProductDetails
}

type OrderProducts struct {
	ProductId string `json:"id"`
	Stock     int    `json:"stock"`
}

type CombinedOrderDetails struct {
	OrderId        string  `json:"order_id"`
	FinalPrice     float64 `json:"final_price"`
	ShipmentStatus string  `json:"shipment_status"`
	PaymentStatus  string  `json:"payment_status"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	Phone          string  `json:"phone"`
	HouseName      string  `json:"house_name" validate:"required"`
	Street         string  `json:"street"`
	City           string  `json:"city"`
	State          string  `json:"state" validate:"required"`
	Pin            string  `json:"pin" validate:"required"`
}

type Page struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type OrderDetailsAdmin struct {
	TotalAmount float64 `gorm:"column:total_amount"`
	ProductName string  `gorm:"column:product_name"`
}
type ItemDetails struct {
	ProductName string  `json:"product_name"`
	FinalPrice  float64 `json:"final_price"`
	Price       float64 `json:"price" `
	Total       float64 `json:"total_price"`
	Quantity    int     `json:"quantity"`
}
