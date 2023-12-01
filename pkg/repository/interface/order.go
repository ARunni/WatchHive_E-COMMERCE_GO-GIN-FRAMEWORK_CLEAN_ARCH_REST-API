package interfaces

import "WatchHive/pkg/utils/models"

type OrderRepository interface {
	GetAllPaymentOption() ([]models.PaymentDetails, error)
	GetAddressFromOrderId(orderID int) (models.AddressInfoResponse, error)
	GetOrderDetailOfAproduct(orderID int) (models.OrderDetails, error)
	GetProductsInCart(cart_id int) ([]int, error)
	FindProductNames(product_id int) (string, error)
	FindCartQuantity(cart_id, product_id int) (int, error)
	FindPrice(product_id int) (float64, error)
	FindStock(id int) (int, error)
}
