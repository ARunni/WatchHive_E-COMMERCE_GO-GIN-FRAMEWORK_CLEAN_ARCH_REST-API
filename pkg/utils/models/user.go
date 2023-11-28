package models

type UserDetails struct {
	Name            string `json:"name"`
	Email           string `json:"email" validate:"email"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmpassword"`
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
	State     string `json:"state" validate:"required"`
	Pin       string `json:"pin" validate:"required"`
}
type UsersProfileDetails struct {
	Firstname    string `json:"firstname" `
	Lastname     string `json:"lastname" `
	Email        string `json:"email" `
	Phone        string `json:"phone" `
	ReferralCode string `json:"referral_code"`
}
