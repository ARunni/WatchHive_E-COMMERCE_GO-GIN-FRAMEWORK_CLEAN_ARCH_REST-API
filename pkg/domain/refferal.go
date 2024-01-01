package domain

import "gorm.io/gorm"

// type Refferal struct {
// 	ID     uint  `json:"id" gorm:"primaryKey"`
// 	UserId uint  `json:"user_id"`
// 	Users  Users `json:"users" gorm:"foreignKey:User_Id"`
// 	RefferalCode string `json:"refferal_code"`

// }

type Referral struct {
	gorm.Model
	UserID         uint    `json:"user_id" gorm:"uniquekey; not null"`
	Users          Users   `json:"-" gorm:"foreignkey:UserID"`
	ReferralCode   string  `json:"referral_code"`
	ReferralAmount float64 `json:"referral_amount"`
	ReferredUserID uint    `json:"referred_user_id"`
}
