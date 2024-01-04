package domain

type Users struct {
	ID       uint   `json:"id" gorm:"unique; not null"`
	Name     string `json:"name"`
	Email    string `json:"email" validate:"email" gorm:"unique; not null"`
	Password string `json:"password" validate:"min=8,max=20"`
	Phone    string `json:"phone" validate:"e164" gorm:"unique; not null"`
	Blocked  bool   `json:"blocked" gorm:"default:false"`
	IsAdmin  bool   `json:"is_admin" gorm:"default:false"`
}

type Address struct {
	Id        uint   `json:"id" gorm:"unique;not null"`
	UserID    uint   `json:"user_id"`
	Users     Users  `json:"-" gorm:"foreignkey:UserID"`
	Name      string `json:"name" validate:"required"`
	HouseName string `json:"house_name" validate:"required"`
	Street    string `json:"street" validate:"required"`
	City      string `json:"city" validate:"required"`
	State     string `json:"state" validate:"required"`
	Phone     string `json:"phone" gorm:"phone,unique"`
	Pin       int    `json:"pin" validate:"required"`
}
