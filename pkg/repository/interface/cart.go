package interfaces

import "WatchHive/pkg/utils/models"

type CartRepository interface {
	AddToCart(userId, productId, Quantity int, productprice float64) error
	CheckStock(product_id int) (int, error)
	QuantityOfProductInCart(userId, productId int) (int, error)
	UpdateCart(quantity int, price float64, userID int, product_id int) error
	TotalPriceForProductInCart(userID, productID int) (float64, error)
	DisplayCart(userID int) ([]models.Cart, error)
	GetTotalPrice(userID int) (models.CartTotal, error)
	UpdateProductQuantityCart(cart models.AddCart) error
	RemoveFromCart(cart models.RemoveFromCart) error
	CheckCart(userID int) (bool, error)

	TotalAmountInCart(userID int) (float64, error)
	UpdateCartAfterOrder(userID, productID int, quantity float64) error
}

