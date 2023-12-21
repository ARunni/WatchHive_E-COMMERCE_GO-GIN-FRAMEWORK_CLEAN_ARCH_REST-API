package usecase

import (
	walletrep "WatchHive/pkg/repository/interface"
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
	"errors"
)

type walletUsecase struct {
	walletRepo walletrep.WalletRepository
}

func NewWalletUsecase(walletRep walletrep.WalletRepository) interfaces.WalletUsecase {
	return &walletUsecase{walletRepo: walletRep}
}

func (wu *walletUsecase) GetWallet(userID int) (models.WalletAmount, error) {
	ok, err := wu.walletRepo.IsWalletExist(userID)
	if err != nil {
		return models.WalletAmount{}, errors.New("error in db")
	}
	if !ok {
		err = wu.walletRepo.CreateWallet(userID)
		if err != nil {
			return models.WalletAmount{}, err
		}
	}
	return wu.walletRepo.GetWallet(userID)
}
