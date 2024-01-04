package domain

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID          int           `json:"user_id" gorm:"not null"`
	Users           Users         `json:"-" gorm:"foreignkey:UserID"`
	AddressID       uint          `json:"address_id" gorm:"not null"`
	Address         Address       `json:"-" gorm:"foreignkey:AddressID"`
	PaymentMethodID uint          `json:"paymentmethod_id"`
	PaymentMethod   PaymentMethod `json:"-" gorm:"foreignkey:PaymentMethodID"`
	ShipmentStatus  string        `json:"shipment_status" gorm:"default:'pending'"`
	PaymentStatus   string        `json:"payment_status" gorm:"default:'not paid'"`
	TotalAmount     float64       `json:"total_amount"`
	FinalPrice      float64       `json:"final_price"`
	Approval        bool          `json:"approval" gorm:"default:false"`
}

type OrderItem struct {
	ID         uint    `json:"id" gorm:"primaryKey;not null"`
	OrderID    uint    `json:"order_id"`
	Order      Order   `json:"-" gorm:"foreignkey:OrderID;constraint:OnDelete:CASCADE"`
	ProductID  uint    `json:"product_id"`
	Products   Product `json:"-" gorm:"foreignkey:ProductID"`
	Quantity   float64 `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}

// type Order struct {
//     gorm.Model
//     UserID          uint          `json:"user_id" gorm:"not null"`
//     Users           Users         `json:"-" gorm:"foreignkey:UserID"`
//     AddressID       uint          `json:"address_id" gorm:"not null"`
//     Address         Address       `json:"-" gorm:"foreignkey:AddressID"`
//     PaymentMethodID uint          `json:"paymentmethod_id"`
//     PaymentMethod   PaymentMethod `json:"-" gorm:"foreignkey:PaymentMethodID"`
//     FinalPrice      float64       `json:"price"`
//     OrderStatus     string        `json:"order_status" gorm:"order_status:4;default:'PENDING';check:order_status IN ('PENDING', 'SHIPPED','DELIVERED','CANCELED','RETURNED')"`
//     PaymentStatus   string        `json:"payment_status" gorm:"payment_status:4;default:'NOT PAID';check:payment_status IN ('PAID', 'NOT PAID','REFUND IN PROGRESS','RETURNED TO WALLET')"`
// }
