package models
type WalletAmount struct {
	Amount float64 `json:"amount"`
}

type WalletHistory struct {
	ID          int     `json:"id"  gorm:"unique;not null"`
	OrderID     int     `json:"order_id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	IsCredited  bool    `json:"is_credited"`
}
