package repository

import (
	"WatchHive/pkg/domain"
	interfaces "WatchHive/pkg/repository/interface"
	"WatchHive/pkg/utils/models"
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type adminRepository struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository {
	return &adminRepository{
		DB: DB,
	}
}

func (ad *adminRepository) LoginHandler(adminDetails models.AdminLogin) (models.Admin, error) {

	var adminCompareDetails models.Admin
	if err := ad.DB.Raw("SELECT * FROM users WHERE email = ? ", adminDetails.Email).Scan(&adminCompareDetails).Error; err != nil {
		return models.Admin{}, err
	}

	return adminCompareDetails, nil
}

func (ad *adminRepository) GetUserByID(id string) (domain.Users, error) {
	user_id, err := strconv.Atoi(id)
	if err != nil {
		return domain.Users{}, err
	}

	var count int64
	if err := ad.DB.Model(&domain.Users{}).Where("id = ?", user_id).Count(&count).Error; err != nil {
		return domain.Users{}, err
	}
	if count < 1 {
		return domain.Users{}, errors.New("user for the given id does not exist")
	}

	var userDetails domain.Users
	if err := ad.DB.Where("id = ?", user_id).First(&userDetails).Error; err != nil {
		return domain.Users{}, err
	}

	return userDetails, nil
}

func (ad *adminRepository) UpdateBlockUserByID(user domain.Users) error {

	err := ad.DB.Exec("update users set blocked = ? where id = ? AND is_admin = 'f'", user.Blocked, user.ID).Error
	if err != nil {
		return err
	}
	return nil

}

func (ad *adminRepository) GetUsers(page int) ([]models.UserDetailsAtAdmin, error) {

	// pagination purpose -
	if page == 0 {
		page = 1
	}

	offset := (page - 1) * 2
	var userDetails []models.UserDetailsAtAdmin

	if err := ad.DB.Raw("select id,name,email,phone,blocked from users  where is_admin = 'f' limit ? offset ?", 3, offset).Scan(&userDetails).Error; err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}
	return userDetails, nil

}

func (ad *adminRepository) IsUserExist(userID int) bool {
	var count int
	err := ad.DB.Raw("select count(*) from users where id = ?", userID).Scan(&count).Error
	if err != nil {
		return false
	}
	if count <= 0 {
		return false
	}
	return true
}

//DashBoard

func (ar *adminRepository) DashboardUserDetails() (models.DashBoardUser, error) {
	var userDetails models.DashBoardUser
	err := ar.DB.Raw("select count(*) from users where is_admin = 'false'").Scan(&userDetails.TotalUsers).Error
	if err != nil {
		err = errors.New("cannot get total users from db")
		return models.DashBoardUser{}, err
	}
	err = ar.DB.Raw("select count(*) from users where blocked = 'true'").Scan(&userDetails.BlockedUser).Error
	if err != nil {
		err = errors.New("cannot get blocked users from db")
		return models.DashBoardUser{}, err
	}
	return userDetails, nil
}

func (ar *adminRepository) DashboardProductDetails() (models.DashBoardProduct, error) {
	var productDetails models.DashBoardProduct
	err := ar.DB.Raw("select count(*) from products").Scan(&productDetails.TotalProducts).Error
	if err != nil {
		err = errors.New("cannot get total products from db")
		return models.DashBoardProduct{}, err
	}
	err = ar.DB.Raw("select count(*) from products where stock <= 0").Scan(&productDetails.OutofStockProduct).Error
	if err != nil {
		err = errors.New("cannot get stock from db")
		return models.DashBoardProduct{}, err
	}
	return productDetails, nil
}

func (ar *adminRepository) DashboardAmountDetails() (models.DashBoardAmount, error) {
	var amountDetails models.DashBoardAmount
	querry :=
		`select coalesce(sum(final_price),0) from orders where payment_status = 'PAID'
	`
	err := ar.DB.Raw(querry).Scan(&amountDetails.CreditedAmount).Error
	if err != nil {
		err = errors.New("cannot get total amount from db")
		return models.DashBoardAmount{}, err
	}
	querry =
		`select coalesce(sum(final_price),0) 
	from orders where payment_status = 'not_paid' 
	and 
	shipment_status = 'pending'
	or 
	shipment_status = 'processing'
	or 
	shipment_status = 'shipped'
	 
	`
	err = ar.DB.Raw(querry).Scan(&amountDetails.PendingAmount).Error
	if err != nil {
		err = errors.New("cannot get pending amount from db")
		return models.DashBoardAmount{}, err
	}
	return amountDetails, nil
}

func (ar *adminRepository) DashboardOrderDetails() (models.DashBoardOrder, error) {
	var orderDetails models.DashBoardOrder
	err := ar.DB.Raw("select count(*) from orders where payment_status = 'PAID'").Scan(&orderDetails.CompletedOrder).Error
	if err != nil {
		err = errors.New("cannot get total order from db")
		return models.DashBoardOrder{}, err
	}
	err = ar.DB.Raw("select count(*) from orders where shipment_status = 'pending' or shipment_status = 'processing'").Scan(&orderDetails.PendingOrder).Error
	if err != nil {
		err = errors.New("cannot get pending orders from db")
		return models.DashBoardOrder{}, err
	}
	err = ar.DB.Raw("select count(*) from orders where shipment_status = 'cancelled' ").Scan(&orderDetails.CancelledOrder).Error
	if err != nil {
		err = errors.New("cannot get cancelled orders from db")
		return models.DashBoardOrder{}, err
	}
	err = ar.DB.Raw("select count(*) from orders  ").Scan(&orderDetails.TotalOrder).Error
	if err != nil {
		err = errors.New("cannot get total orders from db")
		return models.DashBoardOrder{}, err
	}
	err = ar.DB.Raw("select sum(quantity) from order_items  ").Scan(&orderDetails.TotalOrderItem).Error
	if err != nil {
		err = errors.New("cannot get total order items from db")
		return models.DashBoardOrder{}, err
	}
	return orderDetails, nil
}

func (ar *adminRepository) DashboardTotalRevenueDetails() (models.DashBoardRevenue, error) {
	var revenueDetails models.DashBoardRevenue
	startTime := time.Now().AddDate(0, 0, -1)
	err := ar.DB.Raw("select coalesce(sum(final_price),0) from orders where payment_status ='PAID' and created_at >= ?", startTime).Scan(&revenueDetails.TodayRevenue).Error
	if err != nil {
		err = errors.New("cannot get today revenue from db")
		return models.DashBoardRevenue{}, err
	}
	startTime = time.Now().AddDate(0, -1, 1)
	err = ar.DB.Raw("select coalesce(sum(final_price),0) from orders where payment_status ='PAID' and created_at >= ?", startTime).Scan(&revenueDetails.MonthRevenue).Error
	if err != nil {
		err = errors.New("cannot get month revenue from db")
		return models.DashBoardRevenue{}, err
	}
	startTime = time.Now().AddDate(-1, 1, 1)
	err = ar.DB.Raw("select coalesce(sum(final_price),0) from orders where payment_status ='PAID' and created_at >= ?", startTime).Scan(&revenueDetails.YearRevenue).Error
	if err != nil {
		err = errors.New("cannot get year revenue from db")
		return models.DashBoardRevenue{}, err
	}
	return revenueDetails, nil
}

//sales report

func (ar *adminRepository) FilteredSalesReport(startTime time.Time, endTime time.Time) (models.SalesReport, error) {
	var salesReport models.SalesReport
	querry := `
		SELECT COALESCE(SUM(final_price),0) 
		FROM orders WHERE payment_status='PAID'
		AND created_at >= ? AND created_at <= ?
		`
	result := ar.DB.Raw(querry, startTime, endTime).Scan(&salesReport.TotalSales)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}

	result = ar.DB.Raw("SELECT COUNT(*) FROM orders where created_at >= ? AND created_at <= ?", startTime, endTime).Scan(&salesReport.TotalOrders)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}

	querry = `
		SELECT COUNT(*) FROM orders 
		WHERE payment_status = 'PAID' and 
		created_at >= ? AND created_at <= ?
		`

	result = ar.DB.Raw(querry, startTime, endTime).Scan(&salesReport.CompletedOrders)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}

	querry = `
		SELECT COUNT(*) FROM orders WHERE 
		shipment_status = 'processing' AND 
		approval = false AND created_at >= ? AND created_at<=?
		`
	result = ar.DB.Raw(querry, startTime, endTime).Scan(&salesReport.PendingOrders)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	querry = `
		SELECT COUNT(*) FROM orders WHERE 
		shipment_status = 'cancelled' AND created_at >= ? AND created_at<=?
		`
	result = ar.DB.Raw(querry, startTime, endTime).Scan(&salesReport.CancelledOrders)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	querry = `
		SELECT COUNT(*) FROM orders WHERE 
		shipment_status = 'returned' AND created_at >= ? AND created_at<=?
		`
	result = ar.DB.Raw(querry, startTime, endTime).Scan(&salesReport.ReturnedOrders)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}

	var productID int
	querry = `
		SELECT product_id FROM order_items 
		GROUP BY product_id order by SUM(quantity) DESC LIMIT 1
		`
	result = ar.DB.Raw(querry).Scan(&productID)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}

	result = ar.DB.Raw("SELECT product_name FROM products WHERE id = ?", productID).Scan(&salesReport.TrendingProduct)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	return salesReport, nil
}

