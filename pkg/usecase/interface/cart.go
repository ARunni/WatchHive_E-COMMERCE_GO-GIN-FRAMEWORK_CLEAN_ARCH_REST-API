package interfaces

import "WatchHive/pkg/utils/models"

type CartUseCase interface {
	AddToCart(productID int , userID int) (models.CartResponse,error)
}