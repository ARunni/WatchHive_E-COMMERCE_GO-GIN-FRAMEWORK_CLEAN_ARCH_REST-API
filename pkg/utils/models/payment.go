package models

type PaymentDetails struct {
	ID           uint   `json:"id"`
	Payment_Name string `json:"payment_name"`
}

type NewPaymentMethod struct {
	PaymentName string `json:"payment_name"`
}

