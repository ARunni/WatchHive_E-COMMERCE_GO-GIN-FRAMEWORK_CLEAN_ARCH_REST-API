package interfaces

import (
	"WatchHive/pkg/domain"
	"WatchHive/pkg/utils/models"
)

type AdminRepository interface {
	LoginHandler(adminDetails models.AdminLogin) (models.Admin, error)
	GetUserByID(id string) (domain.Users, error)
	UpdateBlockUserByID(user domain.Users) error
	GetUsers(page int) ([]models.UserDetailsAtAdmin, error)
	IsUserExist(userID int)bool

	DashboardUserDetails()(models.DashBoardUser,error)
	DashboardProductDetails()(models.DashBoardProduct,error)
	DashboardOrderDetails()(models.DashBoardOrder,error)
	DashboardTotalRevenueDetails()(models.DashBoardRevenue,error)
	DashboardAmountDetails()(models.DashBoardAmount,error)


}
