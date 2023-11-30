package interfaces

import "WatchHive/pkg/utils/models"

type CartUseCase interface {
	AddToCart(cart models.AddCart) (models.CartResponse, error)
	ListCartItems(userID int) (models.CartResponse, error)
	UpdateProductQuantityCart(cart models.AddCart) (models.CartResponse, error)
	RemoveFromCart(cart models.RemoveFromCart) (models.CartResponse, error)
}
