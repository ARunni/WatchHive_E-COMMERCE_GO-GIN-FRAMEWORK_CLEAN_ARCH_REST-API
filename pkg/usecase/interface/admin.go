package interfaces

import (
	"WatchHive/pkg/utils/models"

	"github.com/jung-kurt/gofpdf"
)

type AdminUseCase interface {
	LoginHandler(adminDetails models.AdminLogin) (models.TockenAdmin, error)
	BlockUser(id string) error
	UnBlockUser(id string) error
	GetUsers(page int) ([]models.UserDetailsAtAdmin, error)
	AdminDashboard() (models.CompleteAdminDashboard, error)
	FilteredSalesReport(timePeriod string) (models.SalesReport, error)
	ExecuteSalesReportByDate(startDate, endDate string) (models.SalesReport, error)
	PrintSalesReport(sales []models.OrderDetailsAdmin) (*gofpdf.Fpdf, error)
	SalesByDate(dayInt int, monthInt int, yearInt int) ([]models.OrderDetailsAdmin, error) 
}
