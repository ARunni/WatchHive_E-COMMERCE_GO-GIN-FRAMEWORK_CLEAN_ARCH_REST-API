package models

type WalletAmount struct {
	Amount float64 `json:"amount"`
}

type WalletHistoryResp struct {
	// WalletID int     `json:"wallet_id"  gorm:"not null"`
	OrderID  int     `json:"order_id" gorm:"not null"`
	Amount   float64 `json:"amount" gorm:"not null"`
	Status   string  `json:"status" gorm:"not null"`
}
type WalletHistory struct {
	WalletID int     `json:"wallet_id"  gorm:"not null"`
	OrderID  int     `json:"order_id" gorm:"not null"`
	Amount   float64 `json:"amount" gorm:"not null"`
	Status   string  `json:"status" gorm:"not null"`
}
type Wallet struct {
	Id     uint    `json:"id"`
	UserID int     `json:"user_id"`
	Amount float64 `json:"amount"`
} 
