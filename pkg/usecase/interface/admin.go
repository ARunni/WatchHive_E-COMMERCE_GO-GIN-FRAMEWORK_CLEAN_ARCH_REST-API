package interfaces

import (
	"WatchHive/pkg/utils/models"
)

type AdminUseCase interface {
	LoginHandler(adminDetails models.AdminLogin) (models.TockenAdmin, error)
	BlockUser(id string) error
	UnBlockUser(id string) error
	GetUsers(page int) ([]models.UserDetailsAtAdmin, error)
}
