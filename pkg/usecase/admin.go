package usecase

import (
	"WatchHive/pkg/domain"
	helper "WatchHive/pkg/helper/interface"
	interfaces "WatchHive/pkg/repository/interface"
	usecase "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"

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

func (ad *adminUseCase) LoginHandler(adminDetails models.AdminLogin) (domain.TockenAdmin, error) {

	adminCompareDetails, err := ad.adminRepository.LoginHandler(adminDetails)
	if err != nil {
		return domain.TockenAdmin{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(adminCompareDetails.Password), []byte(adminDetails.Password))

	if err != nil {
		return domain.TockenAdmin{}, err
	}
	var adminDetailsResponse models.AdminDetailsResponse
	err = copier.Copy(&adminDetailsResponse, &adminCompareDetails)
	if err != nil {
		return domain.TockenAdmin{}, err
	}
	access, refresh, err := ad.helper.GenerateTokenAdmin(adminDetailsResponse)
	if err != nil {
		return domain.TockenAdmin{}, err
	}
	return domain.TockenAdmin{
		Admin:        adminDetailsResponse,
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil

}
