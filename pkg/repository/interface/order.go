package interfaces

import (
	"WatchHive/pkg/domain"
	"WatchHive/pkg/utils/models"
)

type OrderRepository interface {
	GetAllPaymentOption() ([]models.PaymentDetails, error)
	GetAddressFromOrderId(orderID int) (models.AddressInfoResponse, error)
	GetOrderDetailOfAproduct(orderID int) (models.OrderDetails, error)
	GetProductsInCart(cart_id int) ([]int, error)
	FindProductNames(product_id int) (string, error)
	FindCartQuantity(cart_id, product_id int) (int, error)
	FindPrice(product_id int) (float64, error)
	FindStock(id int) (int, error)
	CheckOrderID(orderId int) (bool, error)
	OrderItems(ob models.OrderIncoming, price float64) (int, error)
	AddOrderProducts(order_id int, cart []models.Cart) error
	GetBriefOrderDetails(orderID int) (models.OrderSuccessResponse, error)
	OrderExist(orderID int) (bool, error)
	GetShipmentStatus(orderID int) (string, error)
	UpdateOrder(orderID int) error
	GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error)
	UserOrderRelationship(orderID int, userID int) (int, error)
	GetProductDetailsFromOrders(orderID int) ([]models.OrderProducts, error)
	CancelOrders(orderID int) error
	UpdateQuantityOfProduct(orderProducts []models.OrderProducts) error
	GetAllOrdersAdmin(offset, count int) ([]models.CombinedOrderDetails, error)
	ApproveOrder(orderID int) error
	UpdateStockOfProduct(orderProducts []models.OrderProducts) error
	ApproveCodPaid(orderID int) error
	ReturnOrderCod(orderId int) error
	ReturnOrderRazorPay(orderId int) error 
	ApproveCodReturn(orderID int) error
	GetOrder(orderId int) (domain.Order, error)
	GetDetailedOrderThroughId(orderId int) (models.CombinedOrderDetails, error) 
	ApproveRazorPaid(orderID int) error 
	GetPaymentType(orderID int) (int, error) 
	ApproveRazorDelivered(orderID int) error 
	GetPaymentStatus(orderID int) (string,error)
	GetFinalPriceOrder(orderID int) (float64,error)
}
