package models

type AdminLogin struct {
	Email    string `json:"email,omitempty" validate:"required"`
	Password string `json:"password" validate:"min=8,max=20"`
}
type AdminDetailsResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email" validate:"email"`
}

type AdminResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type UserDetailsAtAdmin struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email" validate:"email"`
	Phone   string `json:"phone"`
	Blocked bool   `json:"blocked"`
}

type Admin struct {
	ID       uint   `json:"id" gorm:"unique;not null"`
	Name     string `json:"name" gorm:"validate:required"`
	Email    string `json:"email" gorm:"validate:required"`
	Password string `json:"password" gorm:"validate:required"`
}
type TockenAdmin struct {
	Admin       AdminDetailsResponse
	AccessToken string
	// RefreshToken string
}

type DashBoardUser struct {
	TotalUsers  int `json:"Totaluser"`
	BlockedUser int `json:"Blockuser"`
}
type DashBoardProduct struct {
	TotalProducts     int `json:"Totalproduct"`
	OutofStockProduct int `json:"Outofstock"`
}
type DashBoardOrder struct {
	CompletedOrder int
	PendingOrder   int
	CancelledOrder int
	TotalOrder     int
	TotalOrderItem int
}
type DashBoardRevenue struct {
	TodayRevenue float64
	MonthRevenue float64
	YearRevenue  float64
}
type DashBoardAmount struct {
	CreditedAmount float64
	PendingAmount  float64
}
type CompleteAdminDashboard struct {
	DashboardUser    DashBoardUser
	DashboardProduct DashBoardProduct
	DashboardRevenue DashBoardRevenue
	DashboardOrder   DashBoardOrder
	DashboardAmount  DashBoardAmount
}
// sales report

type SalesReport struct {
	TotalSales      float64
	TotalOrders     int
	CompletedOrders int
	PendingOrders   int
	ReturnedOrders  int
	CancelledOrders int
	TrendingProduct string
}
