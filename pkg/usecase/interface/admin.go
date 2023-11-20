package interfaces

import (
	"WatchHive/pkg/domain"
	"WatchHive/pkg/utils/models"
)

type AdminUseCase interface {
	LoginHandler(adminDetails models.AdminLogin) (domain.TockenAdmin, error)
}
