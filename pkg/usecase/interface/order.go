package interfaces

import "WatchHive/pkg/utils/models"

type OrderUseCase interface {
	Checkout(userID int) (models.CheckoutDetails, error)
	OrderItemsFromCart(orderFromCart models.OrderFromCart, userID int) (models.OrderSuccessResponse, error)
	ExecutePurchaseCOD(orderID int) error
	GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error)
	CancelOrders(orderID int, userId int) error 
}
