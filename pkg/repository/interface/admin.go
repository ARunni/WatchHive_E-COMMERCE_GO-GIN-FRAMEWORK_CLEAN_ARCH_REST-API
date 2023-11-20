package interfaces

import (
	"WatchHive/pkg/domain"
	"WatchHive/pkg/utils/models"
)

type AdminRepository interface {
	LoginHandler(adminDetails models.AdminLogin) (domain.Admin,error)
	
}