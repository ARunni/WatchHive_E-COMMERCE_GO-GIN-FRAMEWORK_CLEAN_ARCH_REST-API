package interfaces

import "WatchHive/pkg/utils/models"

type CartUseCase interface {
	AddToCart(cart models.AddCart) (models.CartResponse, error)
}