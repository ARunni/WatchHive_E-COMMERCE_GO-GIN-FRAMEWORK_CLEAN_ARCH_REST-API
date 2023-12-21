package usecase

import (
	walletrep "WatchHive/pkg/repository/interface"
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
)

type walletUsecase struct {
	walletRepo walletrep.WalletRepository
}

func NewWalletUsecase(walletRep walletrep.WalletRepository) interfaces.WalletUsecase {
	return &walletUsecase{walletRepo: walletRep}
}

func (wu *walletUsecase) GetWallet(userID int) (models.WalletAmount, error) {
	return wu.walletRepo.GetWallet(userID)
}
