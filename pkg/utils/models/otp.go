package models

type OTPdata struct {
	PhoneNumber string `json:"phone,omitempty"`
}
type VerifyData struct {
	PhoneNumber string `json:"phone,omitempty"`
	Code        string `json:"code,omitempty" validate:"required"`
}
