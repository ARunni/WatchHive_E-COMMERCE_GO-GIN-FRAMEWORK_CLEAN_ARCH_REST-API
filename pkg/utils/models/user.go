package models

type UserDetails struct {
	Name            string `json:"name"`
	Email           string `json:"email" validate:"email"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmpassword"`
	ReferralCode    string `json:"referral_code" gorm:"default:''"`
}

type Address struct {
	Id        uint   `json:"id" gorm:"unique;not null"`
	UserID    uint   `json:"user_id"`
	Name      string `json:"name" validate:"required"`
	HouseName string `json:"house_name" validate:"required"`
	Street    string `json:"street" validate:"required"`
	City      string `json:"city" validate:"required"`
	State     string `json:"state" validate:"required"`
	Pin       string `json:"pin" validate:"required"`
}

type UserDetailsResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email" validate:"email"`
	Phone string `json:"phone"`
}

type TokenUsers struct {
	Users UserDetailsResponse
	Token string
}

type UserLogin struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password"`
}
type UserSignInResponse struct {
	Id       uint   `json:"id"`
	UserID   uint   `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email" validate:"email"`
	Phone    string `Json:"phone"`
	Password string `json:"password"`
}

type AddressInfoResponse struct {
	Id        uint   `json:"id" gorm:"unique;not null"`
	Name      string `json:"name" validate:"required"`
	HouseName string `json:"house_name" validate:"required"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state" validate:"required"`
	Phone     string `json:"phone" gorm:"phone,unique"`
	Pin       string `json:"pin" validate:"required"`
}
type AddressInfo struct {
	Name      string `json:"name" validate:"required"`
	HouseName string `json:"house_name" validate:"required"`
	Street    string `json:"street"`
	City      string `json:"city"`
	Phone     string `json:"phone" gorm:"phone,unique"`
	State     string `json:"state" validate:"required"`
	Pin       string `json:"pin" validate:"required"`
}
type UsersProfileDetails struct {
	ID           uint   `json:"id" gorm:"unique;not null"`
	Name         string `json:"name" gorm:"not null"`
	Email        string `json:"email" validate:"email" `
	Phone        string `json:"phone" `
	ReferralCode string `json:"referral_code"`
}

type ChangePassword struct {
	UserID          uint   `json:"id"`
	CurrentPassWord string `json:"current_password"`
	NewPassWord     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UsersProfileDetailsR struct {
	Name  string `json:"name" gorm:"not null"`
	Email string `json:"email" validate:"email" `
	Phone string `json:"phone" `
}
type ChangePasswordR struct {
	CurrentPassWord string `json:"current_password"`
	NewPassWord     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}
