package repository

import (
	interfaces "WatchHive/pkg/repository/interface"
	"WatchHive/pkg/utils/errmsg"
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

func (wr *WalletDB) AddToWalletHistory(wallet models.WalletHistory) error {
	fmt.Println("wallet history")
	fmt.Println("wallet history", wallet)
	querry := `
	insert into wallet_histories 
	(wallet_id,order_id,amount,status)  
	values (?,?,?,?)
	`
	err := wr.Db.Raw(querry, wallet.WalletID, wallet.OrderID, wallet.Amount, wallet.Status).Error
	if err != nil {
		return errors.New(errmsg.ErrWriteDB)
	}
	return nil
}
func (wr *WalletDB) GetWalletData(userID int) (models.Wallet, error) {
	var wallet models.Wallet
	querry := `
	select * from wallets where user_id = ?
	`
	err := wr.Db.Raw(querry, userID).Scan(&wallet).Error
	if err != nil {
		return models.Wallet{}, errors.New(errmsg.ErrGetDB)
	}
	return wallet, nil
}

func (wr *WalletDB) DebitFromWallet(userID int, amount float64) error {
	fmt.Println("wallet repo debit from wallet amount", amount)
	err := wr.Db.Exec("update wallets set amount = amount - ? where user_id = ? ", amount, userID).Error
	if err != nil {
		return errors.New(errmsg.ErrUpdateDB)
	}
	return nil
}
func (wr *WalletDB) GetWalletHistory(walletId int) ([]models.WalletHistoryResp, error) {
	var wallet []models.WalletHistoryResp
	fmt.Println("....................", walletId)
	err := wr.Db.Raw("select * from wallet_histories where wallet_id = ? ", walletId).Scan(&wallet).Error
	if err != nil {
		return []models.WalletHistoryResp{}, errors.New(errmsg.ErrGetDB)
	}
	return wallet, nil
}
