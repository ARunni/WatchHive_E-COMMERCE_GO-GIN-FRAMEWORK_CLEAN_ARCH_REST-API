package domain

type PaymentMethod struct {
	ID           uint   `json:"id" gorm:"primarykey;not null"`
	Payment_Name string `json:"payment_name" gorm:"unique; not null"`
}
