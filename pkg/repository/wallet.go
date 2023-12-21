package repository

import (
	interfaces "WatchHive/pkg/repository/interface"
	"WatchHive/pkg/utils/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type WalletDB struct {
	Db *gorm.DB
}

func NewWalletRepository(DB *gorm.DB) interfaces.WalletRepository {
	return &WalletDB{Db: DB}
}

func (wr *WalletDB) CreateWallet(userID int) error {

	result := wr.Db.Exec("INSERT INTO wallets (user_id) VALUES (?) RETURNNING id", userID)
	if result.Error != nil {
		fmt.Println("error walletcreation id:")
		return errors.New("cannot create wallet error at database")
	}
	// var walletID int
	// if row := result.RowsAffected; row > 0 {
	// 	err := result.Scan(&walletID).Error
	// 	if err != nil {
	// 		return  errors.New("getting data failed error at database")
	// 	}
	// }

	return nil
}

func (w *WalletDB) GetWallet(userID int) (models.WalletAmount, error) {
	var walletAmount models.WalletAmount
	if err := w.Db.Raw("select amount from wallets where user_id = ?", userID).Scan(&walletAmount).Error; err != nil {
		return models.WalletAmount{}, err
	}
	return walletAmount, nil
}

func (w *WalletDB) IsWalletExist(userID int) (bool, error) {
	var count int
	err := w.Db.Raw("select count(*) from wallets where user_id=?", userID).Scan(&count).Error
	if err != nil {
		return false, errors.New("cannot get wallet details")
	}
	return count >= 1, nil
}
