package interfaces

import (
	"WatchHive/pkg/domain"
	"WatchHive/pkg/utils/models"
	"time"
)

type AdminRepository interface {
	LoginHandler(adminDetails models.AdminLogin) (models.Admin, error)
	GetUserByID(id string) (domain.Users, error)
	UpdateBlockUserByID(user domain.Users) error
	GetUsers(page int) ([]models.UserDetailsAtAdmin, error)
	IsUserExist(userID int) bool

	DashboardUserDetails() (models.DashBoardUser, error)
	DashboardProductDetails() (models.DashBoardProduct, error)
	DashboardOrderDetails() (models.DashBoardOrder, error)
	DashboardTotalRevenueDetails() (models.DashBoardRevenue, error)
	DashboardAmountDetails() (models.DashBoardAmount, error)

	FilteredSalesReport(startTime time.Time, endTime time.Time) (models.SalesReport, error)

	SalesByDay(yearInt int, monthInt int, dayInt int) ([]models.OrderDetailsAdmin, error)
	SalesByMonth(yearInt int, monthInt int) ([]models.OrderDetailsAdmin, error)
	SalesByYear(yearInt int) ([]models.OrderDetailsAdmin, error)
}
