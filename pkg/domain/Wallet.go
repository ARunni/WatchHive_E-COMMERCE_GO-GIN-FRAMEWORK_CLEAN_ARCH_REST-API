package domain

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	Id     uint    `json:"id" gorm:"unique;not null" `
	UserID int     `json:"user_id"`
	Users  Users   `json:"-" gorm:"foreignkey:UserID"`
	Amount float64 `json:"amount" gorm:"default:0.0"`
}

type WalletHistory struct {
	ID       uint    `json:"id"  gorm:"unique;not null"`
	WalletID int     `json:"wallet_id" gorm:"not null"`
	Wallet   Wallet  `json:"-" gorm:"foreignkey:WalletID"`
	OrderID  int     `json:"order_id" gorm:"not null"`
	Amount   float64 `json:"amount" gorm:"not null"`
	Status   string  `json:"status" gorm:"status:2;default:'CREDITED';check:status IN ('CREDITED','DEBITED')"`
}
