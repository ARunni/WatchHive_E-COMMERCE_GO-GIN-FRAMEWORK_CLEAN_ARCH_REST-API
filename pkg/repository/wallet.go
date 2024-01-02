package repository

import (
	interfaces "WatchHive/pkg/repository/interface"
	"WatchHive/pkg/utils/errmsg"
	"WatchHive/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)

type WalletDB struct {
	Db *gorm.DB
}

func NewWalletRepository(DB *gorm.DB) interfaces.WalletRepository {
	return &WalletDB{Db: DB}
}

func (wr *WalletDB) CreateWallet(userID int) error {

	err := wr.Db.Exec("INSERT INTO wallets (created_at ,user_id) VALUES (NOW(),?) RETURNING id", userID).Error
	if err != nil {

		return errors.New(errmsg.ErrWriteDB)
	}

	return nil
}

func (wr *WalletDB) GetWallet(userID int) (models.WalletAmount, error) {
	var walletAmount models.WalletAmount
	if err := wr.Db.Raw("select amount from wallets where user_id = ?", userID).Scan(&walletAmount).Error; err != nil {
		return models.WalletAmount{}, err
	}
	return walletAmount, nil
}

func (wr *WalletDB) IsWalletExist(userID int) (bool, error) {
	var count int
	err := wr.Db.Raw("select count(*) from wallets where user_id=?", userID).Scan(&count).Error
	if err != nil {
		return false, errors.New(errmsg.ErrGetDB)
	}
	return count >= 1, nil
}

func (wr *WalletDB) AddToWallet(userID int, Amount float64) error {
	err := wr.Db.Exec("update wallets set amount = amount+? where user_id = ? returning amount", Amount, userID).Error
	if err != nil {
		return errors.New(errmsg.ErrWriteDB)
	}
	return nil

}
