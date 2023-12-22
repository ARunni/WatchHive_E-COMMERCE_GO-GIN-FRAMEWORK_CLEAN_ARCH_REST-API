package interfaces

import "WatchHive/pkg/utils/models"

type WalletRepository interface {
	CreateWallet(userID int) error
	GetWallet(userID int) (models.WalletAmount, error)
	IsWalletExist(userID int) (bool, error) 
	AddToWallet(userID int,Amount float64) error
}
