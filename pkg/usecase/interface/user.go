package interfaces

import "WatchHive/pkg/utils/models"

type UserUseCase interface {
	UserSignUp(user models.UserDetails) (models.TokenUsers, error)
	LoginHandler(user models.UserLogin) (models.TokenUsers, error)
	AddAddress(userID int, address models.AddressInfoResponse) (models.AddressInfoResponse,error)

}
