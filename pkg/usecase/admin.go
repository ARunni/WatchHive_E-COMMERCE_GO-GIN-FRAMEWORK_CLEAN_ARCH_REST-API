package usecase

import (
	helper "WatchHive/pkg/helper/interface"
	interfaces "WatchHive/pkg/repository/interface"
	usecase "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
	"errors"
	"strconv"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type adminUseCase struct {
	adminRepository interfaces.AdminRepository
	helper          helper.Helper
}

func NewAdminUseCase(repo interfaces.AdminRepository, h helper.Helper) usecase.AdminUseCase {
	return &adminUseCase{
		adminRepository: repo,
		helper:          h,
	}
}

func (ad *adminUseCase) LoginHandler(adminDetails models.AdminLogin) (models.TockenAdmin, error) {

	adminCompareDetails, err := ad.adminRepository.LoginHandler(adminDetails)
	if err != nil {
		return models.TockenAdmin{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(adminCompareDetails.Password), []byte(adminDetails.Password))
	if err != nil {
		return models.TockenAdmin{}, err
	}
	var adminDetailsResponse models.AdminDetailsResponse
	err = copier.Copy(&adminDetailsResponse, &adminCompareDetails)
	if err != nil {
		return models.TockenAdmin{}, err
	}
	access, _, err := ad.helper.GenerateTokenAdmin(adminDetailsResponse)
	if err != nil {
		return models.TockenAdmin{}, err
	}
	return models.TockenAdmin{
		Admin:       adminDetailsResponse,
		AccessToken: access,
		// RefreshToken: refresh,
	}, nil

}

func (ad *adminUseCase) BlockUser(id string) error {

	iD, IdErr := strconv.Atoi(id)
	if IdErr != nil || iD <= 0 {
		return errors.New("invalid id")
	}

	user, err := ad.adminRepository.GetUserByID(id)
	if err != nil {
		return err
	}

	if user.Blocked {
		return errors.New("already blocked")
	} else {
		user.Blocked = true
	}

	err = ad.adminRepository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}

	return nil

}

func (ad *adminUseCase) UnBlockUser(id string) error {

	iD, IdErr := strconv.Atoi(id)
	if IdErr != nil || iD <= 0 {
		return errors.New("invalid id")
	}

	user, err := ad.adminRepository.GetUserByID(id)
	if err != nil {
		return err
	}

	if user.Blocked {
		user.Blocked = false
	} else {
		return errors.New("already unblocked")
	}

	err = ad.adminRepository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}

	return nil

}
func (ad *adminUseCase) GetUsers(page int) ([]models.UserDetailsAtAdmin, error) {
	if page < 0 {
		return []models.UserDetailsAtAdmin{}, errors.New("page number cannot be negative")
	}

	userDetails, err := ad.adminRepository.GetUsers(page)
	if err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}

	return userDetails, nil

}

///1 DEC
