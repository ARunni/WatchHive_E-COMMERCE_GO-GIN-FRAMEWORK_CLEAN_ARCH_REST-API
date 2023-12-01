package interfaces

import "WatchHive/pkg/utils/models"






type OrderUseCase interface {
	Checkout(userID int) (models.CheckoutDetails, error)
}
