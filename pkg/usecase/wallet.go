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

	amount, err := wu.walletRepo.GetWallet(userID)
	if err != nil {
		return models.WalletAmount{}, err
	}
	return amount, nil
}

func (wu *walletUsecase) GetWalletHistory(userId int) ([]models.WalletHistoryResp, error) {

	wallet, err := wu.walletRepo.GetWalletData(userId)
	if err != nil {
		return []models.WalletHistoryResp{}, err
	}
	walletResp, err := wu.walletRepo.GetWalletHistory(int(wallet.ID))
	if err != nil {
		return []models.WalletHistoryResp{}, err
	}

	return walletResp, nil
}
