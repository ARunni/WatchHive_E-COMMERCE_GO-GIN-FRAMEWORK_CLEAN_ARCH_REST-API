package domain

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	Id     uint    `json:"id" gorm:"unique;not null" `
	UserID int     `json:"user_id"`
	Users  Users   `json:"-" gorm:"foreignkey:UserID"`
	Amount float64 `json:"amount" gorm:"default:0.0"`
}
