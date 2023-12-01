package interfaces

import (
	"WatchHive/pkg/domain"
	"WatchHive/pkg/utils/models"
)

type AdminRepository interface {
	LoginHandler(adminDetails models.AdminLogin) (domain.Admin, error)
	GetUserByID(id string) (domain.Users, error)
	UpdateBlockUserByID(user domain.Users) error
	GetUsers(page int) ([]models.UserDetailsAtAdmin, error)
	CheckIfPaymentMethodAlreadyExists(payment string) (bool, error)
	AddPaymentMethod(pay models.NewPaymentMethod) (domain.PaymentMethod, error)
}
