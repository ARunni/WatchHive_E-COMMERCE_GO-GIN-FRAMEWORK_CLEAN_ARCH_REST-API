package usecase

import (
	helper "WatchHive/pkg/helper/interface"
	interfaces "WatchHive/pkg/repository/interface"
	usecase "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
	"errors"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	"github.com/jung-kurt/gofpdf"
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
	ok := ad.adminRepository.IsUserExist(iD)
	if !ok {
		return errors.New("user does not exist")
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
	ok := ad.adminRepository.IsUserExist(iD)
	if !ok {
		return errors.New("user does not exist")
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

func (au *adminUseCase) AdminDashboard() (models.CompleteAdminDashboard, error) {
	userDetails, err := au.adminRepository.DashboardUserDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	productDetails, err := au.adminRepository.DashboardProductDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	orderDetails, err := au.adminRepository.DashboardOrderDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	amountDetails, err := au.adminRepository.DashboardAmountDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	totalRevnueDetails, err := au.adminRepository.DashboardTotalRevenueDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	return models.CompleteAdminDashboard{
		DashboardUser:    userDetails,
		DashboardProduct: productDetails,
		DashboardOrder:   orderDetails,
		DashboardAmount:  amountDetails,
		DashboardRevenue: totalRevnueDetails,
	}, nil
}

// sales report
func (ah *adminUseCase) FilteredSalesReport(timePeriod string) (models.SalesReport, error) {
	if timePeriod == "" {
		err := errors.New("field cannot be empty")
		return models.SalesReport{}, err
	}

	if timePeriod != "week" && timePeriod != "month" && timePeriod != "year" {
		err := errors.New("invalid time period, available options : week, month & year")
		return models.SalesReport{}, err
	}

	startTime, endTime := ah.helper.GetTimeFromPeriod(timePeriod)
	saleReport, err := ah.adminRepository.FilteredSalesReport(startTime, endTime)

	if err != nil {
		return models.SalesReport{}, err
	}
	return saleReport, nil
}

func (au *adminUseCase) ExecuteSalesReportByDate(startDate, endDate string) (models.SalesReport, error) {

	parsedStartDate, err := time.Parse("02-01-2006", startDate)
	if err != nil {
		err := errors.New("enter the date in correct format")
		return models.SalesReport{}, err
	}

	isValid := !parsedStartDate.IsZero()
	if !isValid {
		err := errors.New("enter the date in correct format & valid date")
		return models.SalesReport{}, err
	}
	parsedEndDate, err := time.Parse("02-01-2006", endDate)
	if err != nil {
		err := errors.New("enter the date in correct format")
		return models.SalesReport{}, err
	}

	isValid = !parsedEndDate.IsZero()
	if !isValid {
		err := errors.New("enter the date in correct format & valid date")
		return models.SalesReport{}, err
	}

	if parsedStartDate.After(parsedEndDate) {
		err := errors.New("start date is after end date")

		return models.SalesReport{}, err
	}

	orders, err := au.adminRepository.FilteredSalesReport(parsedStartDate, parsedEndDate)
	if err != nil {
		return models.SalesReport{}, errors.New("report fetching failed")
	}
	return orders, nil
}
func (ad *adminUseCase) SalesByDate(dayInt int, monthInt int, yearInt int) ([]models.OrderDetailsAdmin, error) {

	if dayInt == 0 && monthInt == 0 && yearInt == 0 {
		return []models.OrderDetailsAdmin{}, errors.New("must enter a value for day, month, and year")
	}

	if dayInt < 0 || monthInt < 0 || yearInt < 0 {
		return []models.OrderDetailsAdmin{}, errors.New("no such values are allowded")
	}

	if yearInt >= 2020 {
		if monthInt == 0 && dayInt == 0 {

			body, err := ad.adminRepository.SalesByYear(yearInt)
			if err != nil {
				return []models.OrderDetailsAdmin{}, err
			}

			return body, nil
		} else if monthInt > 0 && monthInt <= 12 && dayInt == 0 {

			body, err := ad.adminRepository.SalesByMonth(yearInt, monthInt)
			if err != nil {
				return []models.OrderDetailsAdmin{}, err
			}

			return body, nil
		} else if monthInt > 0 && monthInt <= 12 && dayInt > 0 && dayInt <= 31 {

			body, err := ad.adminRepository.SalesByDay(yearInt, monthInt, dayInt)
			if err != nil {
				return []models.OrderDetailsAdmin{}, err
			}

			return body, nil
		}
	}

	return []models.OrderDetailsAdmin{}, errors.New("invalid date parameters")
}

func (ad *adminUseCase) PrintSalesReport(sales []models.OrderDetailsAdmin) (*gofpdf.Fpdf, error) {

	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 22)
	pdf.SetTextColor(31, 73, 125)
	pdf.CellFormat(0, 20, "Total Sales Report", "0", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 16)
	pdf.SetTextColor(0, 0, 0)
	var FinalAmount float64
	for _, item := range sales {
		pdf.CellFormat(0, 10, "Product: "+item.ProductName, "0", 1, "L", false, 0, "")
		amount := strconv.FormatFloat(item.TotalAmount, 'f', 2, 64)
		pdf.CellFormat(0, 10, "Amount Sold: $"+amount, "0", 1, "L", false, 0, "")
		pdf.Ln(5)
		FinalAmount += item.TotalAmount
	}
	pdf.SetFont("Arial", "", 18)
	pdf.SetTextColor(0, 0, 0)
	FinalTotal := strconv.FormatFloat(FinalAmount, 'f', 2, 64)
	pdf.CellFormat(0, 10, " Total Amount Sold: "+FinalTotal, "0", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "I", 12)
	pdf.SetTextColor(150, 150, 150)

	pdf.Cell(0, 10, "Generated by Watch Hive India Pvt Ltd. - "+time.Now().Format("2006-01-02 15:04:05"))

	return pdf, nil
}
