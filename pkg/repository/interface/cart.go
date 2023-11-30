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

	// GetTotalPrice(userID int) (models.CartTotal, error)

	// ProductExist(userID int, productID int) (bool, error)
	// CartExist(userID int) (bool, error)
	// EmptyCart(userID int) error
	// RemoveProductFromCart(userID int, product_id int) error

	// QuantityOfProductInCart(userId int, productId int) (int, error)
	// TotalPriceForProductInCart(userID int, productID int) (float64, error)

	// UpdateCart(quantity int, price float64, userID int, product_id int) error

	// GetAllItemsFromCart(userID int) ([]models.Cart, error)
	// GetTotalPriceFromCart(userID int) (float64, error)

	// DisplayCart(userID int) ([]models.Cart, error)
	// CartAfterRemovalOfProduct(user_id int) ([]models.Cart, error)

	// GetQuantityAndProductDetails(userId int, productId int, cartDetails struct {
	// 	Quantity   int
	// 	TotalPrice float64
	// }) (struct {
	// 	Quantity   int
	// 	TotalPrice float64
	// }, error)

	// UpdateCartDetails(cartDetails struct {
	// 	Quantity   int
	// 	TotalPrice float64
	// }, userId int, productId int) error
}
